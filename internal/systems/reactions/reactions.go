/*
 * owobot - The coolest Discord bot ever written
 * Copyright (C) 2023 owobot Contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package reactions

import (
	"fmt"
	"math/rand"
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
	"go.elara.ws/owobot/internal/systems/commands"
	"go.elara.ws/owobot/internal/util"
)

func Init(s *discordgo.Session) error {
	s.AddHandler(onMessage)

	commands.Register(s, reactionsCmd, &discordgo.ApplicationCommand{
		Name:                     "reactions",
		Description:              "Manage message reactions",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageEmojis),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Add a new message reaction",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "match_type",
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "The matcher type for this reaction",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "contains",
								Value: "contains",
							},
							{
								Name:  "regex",
								Value: "regex",
							},
						},
					},
					{
						Name:        "match",
						Description: "What the matcher should look for",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "reaction_type",
						Description: "The reaction type",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "emoji",
								Value: "emoji",
							},
							{
								Name:  "text",
								Value: "text",
							},
						},
					},
					{
						Name:        "reaction",
						Description: "The contents of the reaction",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "chance",
						Description: "The percent chance that the reaction occurs",
						MinValue:    util.Pointer[float64](1),
						MaxValue:    100,
						Type:        discordgo.ApplicationCommandOptionInteger,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "List all the reactions for this guild",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "delete",
				Description: "Remove all message reactions with the given match",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "match",
						Description: "The match value for which to remove reactions",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			},
		},
	})

	return nil
}

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

			var content string
			switch reaction.ReactionType {
			case db.ReactionTypeText:
				submatch := re.FindSubmatch([]byte(mc.Content))
				if len(submatch) > 1 {
					replacements := map[string]any{}
					for i, match := range submatch {
						replacements[strconv.Itoa(i)] = match
					}
					content = fasttemplate.ExecuteStringStd(reaction.Reaction, "{", "}", replacements)
				} else if len(submatch) == 1 {
					content = reaction.Reaction
				}
			case db.ReactionTypeEmoji:
				if re.MatchString(mc.Content) {
					content = reaction.Reaction
				}
			}

			if content != "" {
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

func performReaction(s *discordgo.Session, reaction db.Reaction, content string, mc *discordgo.MessageCreate) error {
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
		_, err := s.ChannelMessageSendReply(mc.ChannelID, content, mc.Reference())
		if err != nil {
			return err
		}
	case db.ReactionTypeEmoji:
		var emojis []string
		if strings.Contains(content, ",") {
			emojis = strings.Split(content, ",")
		} else {
			emojis = []string{content}
		}

		for _, emojiStr := range emojis {
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
