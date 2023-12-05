/*
 * owobot - Your server's guardian and entertainer
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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/emoji"
	"go.elara.ws/owobot/internal/util"
)

func reactionsCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()

	switch name := data.Options[0].Name; name {
	case "add":
		return reactionsAddCmd(s, i)
	case "list":
		return reactionsListCmd(s, i)
	case "delete":
		return reactionsDeleteCmd(s, i)
	default:
		return fmt.Errorf("unknown reactions subcommand: %s", name)
	}
}

func reactionsAddCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	reaction := db.Reaction{
		MatchType:    db.MatchType(args[0].StringValue()),
		Match:        strings.TrimSpace(args[1].StringValue()),
		ReactionType: db.ReactionType(args[2].StringValue()),
		Reaction:     strings.TrimSpace(args[3].StringValue()),
		Chance:       100,
	}

	if len(args) == 5 {
		reaction.Chance = int(args[4].IntValue())
	}

	switch reaction.MatchType {
	case db.MatchTypeRegex:
		if _, err := cache.Regex(reaction.Match); err != nil {
			return err
		}
	case db.MatchTypeContains:
		// Ensure the contains match is lowercase so we can check it
		// against a lowercase string later, in the message handler.
		reaction.Match = strings.ToLower(reaction.Match)
	}

	if reaction.ReactionType == db.ReactionTypeEmoji {
		if err := validateEmoji(reaction.Reaction); err != nil {
			return err
		}
	}

	err := db.AddReaction(i.GuildID, reaction)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, "Successfully added reaction!")
}

func reactionsListCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	reactions, err := db.Reactions(i.GuildID)
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("**Reactions:**\n")
	for _, reaction := range reactions {
		sb.WriteString("- _[")
		if reaction.Chance < 100 {
			sb.WriteString(strconv.Itoa(reaction.Chance))
			sb.WriteString("% ")
		}
		sb.WriteString(string(reaction.MatchType))
		sb.WriteString("]_ `")
		sb.WriteString(reaction.Match)
		sb.WriteString("`: \"")
		sb.WriteString(reaction.Reaction)
		sb.WriteString("\" _(")
		sb.WriteString(string(reaction.ReactionType))
		sb.WriteString(")_\n")
	}

	return util.RespondEphemeral(s, i.Interaction, sb.String())
}

func reactionsDeleteCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Make sure the user has the manage expressions permission
	// in case a role/member override allows someone else to use it
	if i.Member.Permissions&discordgo.PermissionManageEmojis == 0 {
		return errors.New("you do not have permission to delete reactions")
	}

	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	err := db.DeleteReaction(i.GuildID, args[0].StringValue())
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, "Successfully removed reaction")
}

func validateEmoji(s string) error {
	if strings.Contains(s, ",") {
		split := strings.Split(s, ",")
		for _, emojiStr := range split {
			if _, ok := emoji.Parse(emojiStr); !ok {
				return fmt.Errorf("invalid reaction emoji: %s", emojiStr)
			}
		}
	} else {
		if _, ok := emoji.Parse(s); !ok {
			return fmt.Errorf("invalid reaction emoji: %s", s)
		}
	}
	return nil
}
