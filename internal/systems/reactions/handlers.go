package reactions

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/valyala/fasttemplate"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/emoji"
)

// onMessage handles all new messages. It checks if the message matches any reaction
// registered for that guild, and if it does, it performs all the matching reactions.
func onMessage(s *discordgo.Session, mc *discordgo.MessageCreate) {
	if mc.Author.ID == s.State.User.ID {
		return
	}

	reactions, err := db.Reactions(mc.GuildID)
	if err != nil {
		log.Error("Error getting reactions from database").Err(err).Send()
		return
	}

	for _, reaction := range reactions {
		if slices.Contains(reaction.ExcludedChannels, mc.ChannelID) {
			continue
		}

		switch reaction.MatchType {
		case db.MatchTypeContains:
			if strings.Contains(strings.ToLower(mc.Content), reaction.Match) {
				err = performReaction(s, reaction, reaction.Reaction, mc)
				if err != nil {
					log.Error("Error performing reaction").Err(err).Send()
					continue
				}
			}
		case db.MatchTypeRegex:
			re, err := cache.Regex(reaction.Match)
			if err != nil {
				log.Error("Error compiling regex").Err(err).Send()
				continue
			}

			var content db.StringSlice
			switch reaction.ReactionType {
			case db.ReactionTypeText:
				submatch := re.FindSubmatch([]byte(mc.Content))
				if len(submatch) > 1 {
					replacements := map[string]any{}
					for i, match := range submatch {
						replacements[strconv.Itoa(i)] = match
					}
					content = db.StringSlice{
						fasttemplate.ExecuteStringStd(reaction.Reaction[0], "{", "}", replacements),
					}
				} else if len(submatch) == 1 {
					content = reaction.Reaction
				}
			case db.ReactionTypeEmoji:
				if re.MatchString(mc.Content) {
					content = reaction.Reaction
				}
			}

			if content != nil {
				err = performReaction(s, reaction, content, mc)
				if err != nil {
					log.Error("Error performing reaction").Err(err).Send()
					continue
				}
			}
		}
	}
}

var (
	rngMtx = sync.Mutex{}
	rng    = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func performReaction(s *discordgo.Session, reaction db.Reaction, content db.StringSlice, mc *discordgo.MessageCreate) error {
	if reaction.Chance < 100 {
		rngMtx.Lock()
		randNum := rng.Intn(100) + 1
		rngMtx.Unlock()
		if randNum > reaction.Chance {
			return nil
		}
	}

	switch reaction.ReactionType {
	case db.ReactionTypeText:
		_, err := s.ChannelMessageSendReply(mc.ChannelID, content[0], mc.Reference())
		if err != nil {
			return err
		}
	case db.ReactionTypeEmoji:
		for _, emojiStr := range content {
			e, ok := emoji.Parse(emojiStr)
			if !ok {
				return fmt.Errorf("invalid emoji: %s", emojiStr)
			}

			err := s.MessageReactionAdd(mc.ChannelID, mc.ID, e.APIFormat())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
