package polls

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/emoji"
	"go.elara.ws/owobot/internal/util"
	"go.elara.ws/owobot/internal/xsync"
)

// onPollAddOpt handles the Add Option button on unfinished polls.
func onPollAddOpt(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type != discordgo.InteractionMessageComponent {
		return nil
	}

	data := i.MessageComponentData()
	if data.CustomID != "poll-add-opt" {
		return nil
	}

	if i.Member.User.ID != i.Message.Interaction.User.ID {
		return errors.New("only the creator of the poll may add options to it")
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Add an option",
			CustomID: "poll-opt-modal",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						Label:    "Option text",
						Style:    discordgo.TextInputShort,
						Required: true,
					},
				}},
			},
		},
	})
}

// onAddOptModalSubmit handles the submission of the Add Option modal.
func onAddOptModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type != discordgo.InteractionModalSubmit {
		return nil
	}

	data := i.ModalSubmitData()
	if data.CustomID != "poll-opt-modal" {
		return nil
	}

	row := data.Components[0].(*discordgo.ActionsRow)
	textInput := row.Components[0].(*discordgo.TextInput)

	err := db.AddPollOptionText(i.Message.ID, textInput.Value)
	if err != nil {
		return err
	}

	return updatePollWaitReaction(s, i.Interaction)
}

// onPollReaction handles reactions to a poll.
func onPollReaction(s *discordgo.Session, mra *discordgo.MessageReactionAdd) {
	poll, err := db.GetPoll(mra.MessageID)
	if errors.Is(err, sql.ErrNoRows) {
		return
	} else if err != nil {
		log.Error("Error getting poll from database").Err(err).Send()
		return
	}

	err = s.MessageReactionRemove(mra.ChannelID, mra.MessageID, mra.Emoji.APIName(), mra.UserID)
	if err != nil {
		log.Error("Error removing poll reaction").Err(err).Send()
		return
	}

	// If the poll is finished, there's already an emoji for every option,
	// or the user who reacted is not the owner of the poll, return.
	if poll.Finished ||
		len(poll.OptionEmojis) == len(poll.OptionText) ||
		mra.Member.User.ID != poll.OwnerID {
		return
	}

	err = db.AddPollOptionEmoji(mra.MessageID, mra.Emoji.MessageFormat())
	if err != nil {
		log.Error("Error adding poll option emoji").Err(err).Send()
		return
	}

	err = updatePollUnfinished(s, mra.MessageID, mra.ChannelID)
	if err != nil {
		log.Error("Error updating poll message").Err(err).Send()
		return
	}
}

// onPollFinish handles the Finish button on an unfinished poll.
func onPollFinish(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type != discordgo.InteractionMessageComponent {
		return nil
	}

	data := i.MessageComponentData()
	if data.CustomID != "poll-finish" {
		return nil
	}

	if i.Member.User.ID != i.Message.Interaction.User.ID {
		return errors.New("only the creator of the poll may finish it")
	}

	poll, err := db.GetPoll(i.Message.ID)
	if err != nil {
		return err
	}

	privacyToken, err := makePrivacyToken()
	if err != nil {
		return err
	}

	var (
		components []discordgo.MessageComponent
		currentRow discordgo.ActionsRow
	)

	for i := 0; i < len(poll.OptionEmojis); i++ {
		// Action rows can only contain 5 elements,
		// so we create a new row if we reach a multiple
		// of 5.
		if i > 0 && i%5 == 0 {
			components = append(components, currentRow)
			currentRow = discordgo.ActionsRow{}
		}

		e, ok := emoji.Parse(poll.OptionEmojis[i])
		if !ok {
			return fmt.Errorf("invalid emoji: %s", poll.OptionEmojis[i])
		}

		currentRow.Components = append(currentRow.Components, discordgo.Button{
			CustomID: "vote:" + strconv.Itoa(i) + ":" + privacyToken,
			Style:    discordgo.SecondaryButton,
			Emoji: &discordgo.ComponentEmoji{
				Name: e.Name,
				ID:   e.ID,
			},
		})
	}

	if len(currentRow.Components) != 0 {
		components = append(components, currentRow)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    i.Message.Content,
			Components: components,
		},
	})
	if err != nil {
		return err
	}

	err = db.FinishPoll(i.Message.ID)
	if err != nil {
		log.Error("Error finishing poll").Err(err).Send()
		return nil
	}

	_, err = s.MessageThreadStart(i.ChannelID, i.Message.ID, poll.Title, 1440)
	if err != nil {
		log.Error("Error starting poll thread").Err(err).Send()
		return nil
	}

	return nil
}

// voteMtx ensures only one vote is being processed at a time
// for each message, since handlers may be executed concurrently.
var voteMtx xsync.KeyedMutex

// onVote handles poll votes.
func onVote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	data := i.MessageComponentData()
	splitID := strings.SplitN(data.CustomID, ":", 3)
	if splitID[0] != "vote" {
		return
	}

	// Respond with a deferred update since we have no idea how long this could take
	// because we have to get the vote lock which may take a while depending on how many
	// votes are being processed for this poll at the moment.
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Error("Error responding to interaction").Err(err).Send()
	}

	voteMtx.Lock(i.Message.ID)
	defer voteMtx.Unlock(i.Message.ID)

	option, err := strconv.Atoi(splitID[1])
	if err != nil {
		log.Error("Error parsing button option").Err(err).Send()
		return
	}

	err = db.AddVote(db.Vote{
		PollMsgID: i.Message.ID,
		UserToken: makeUserToken(splitID[2], i.Member.User.ID),
		Option:    option,
	})
	if err != nil {
		log.Error("Error adding vote to database").Err(err).Send()
		return
	}

	poll, err := db.GetPoll(i.Message.ID)
	if err != nil {
		log.Error("Error getting poll from database").Err(err).Send()
		return
	}

	content, err := generatePollContent(poll)
	if err != nil {
		log.Error("Error generating poll content").Err(err).Send()
		return
	}

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      i.Message.ID,
		Channel: i.ChannelID,
		Content: util.Pointer(content),
	})
	if err != nil {
		log.Error("Error editing poll message").Err(err).Send()
		return
	}
}

// makePrivacyToken creates a random token to be hashed with the user's
// id to create an anonymous but deterministic user token for use in
// poll votes.
func makePrivacyToken() (string, error) {
	randomData := make([]byte, 25)
	_, err := io.ReadFull(rand.Reader, randomData)
	return base64.URLEncoding.EncodeToString(randomData), err
}

// makeUserToken hashes a privacy token with a user ID to create
// a user token.
func makeUserToken(privacyToken, userID string) string {
	userTokenData := sha256.Sum256(append([]byte(privacyToken), userID...))
	return base64.URLEncoding.EncodeToString(userTokenData[:])
}

// updatePollWaitReaction updates a poll to wait for an emoji reaction.
func updatePollWaitReaction(s *discordgo.Session, i *discordgo.Interaction) error {
	content := i.Message.Content
	content += "\n\n_Please react with the emoji that you'd like to use for this option._"
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
}

// updatePollUnfinished updates an unfinished poll.
func updatePollUnfinished(s *discordgo.Session, msgID, channelID string) error {
	poll, err := db.GetPoll(msgID)
	if err != nil {
		return err
	}

	content, err := generatePollContent(poll)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      msgID,
		Channel: channelID,
		Content: &content,
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Add Option",
					Style:    discordgo.PrimaryButton,
					CustomID: "poll-add-opt",
				},
				discordgo.Button{
					Label:    "Finish",
					Style:    discordgo.SuccessButton,
					CustomID: "poll-finish",
				},
			}},
		},
	})
	return err
}

// generatePollContent generates the body of the poll, with the title and
// option list.
func generatePollContent(poll *db.Poll) (string, error) {
	var sb strings.Builder

	sb.WriteString("**")
	sb.WriteString(poll.Title)
	sb.WriteString("**\n")

	for i, emoji := range poll.OptionEmojis {
		sb.WriteString(emoji)
		sb.WriteByte(' ')
		voteAmount, err := db.VoteAmount(poll.MsgID, i)
		if err != nil {
			return "", err
		}
		sb.WriteString(strconv.Itoa(int(voteAmount)))
		sb.WriteByte(' ')
		sb.WriteString(poll.OptionText[i])
		sb.WriteByte('\n')
	}

	return sb.String(), nil
}
