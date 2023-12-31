package polls

import (
	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/db"
)

func pollCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	title := data.Options[0].StringValue()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "**" + title + "**",
			Components: []discordgo.MessageComponent{
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
		},
	})
	if err != nil {
		return err
	}

	msg, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		return err
	}
	return db.CreatePoll(msg.ID, i.Member.User.ID, title)
}
