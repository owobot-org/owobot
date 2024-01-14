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

// reactionsCmd handles the `/reactions` command and routes it to the correct subcommand.
func reactionsCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()

	switch name := data.Options[0].Name; name {
	case "add":
		return reactionsAddCmd(s, i)
	case "list":
		return reactionsListCmd(s, i)
	case "delete":
		return reactionsDeleteCmd(s, i)
	case "exclude":
		return reactionsExcludeCmd(s, i)
	case "unexclude":
		return reactionsUnexcludeCmd(s, i)
	default:
		return fmt.Errorf("unknown reactions subcommand: %s", name)
	}
}

// reactionsAddCmd handles the `/reactions add` command.
func reactionsAddCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	reaction := db.Reaction{
		MatchType:    db.MatchType(args[0].StringValue()),
		Match:        strings.TrimSpace(args[1].StringValue()),
		ReactionType: db.ReactionType(args[2].StringValue()),
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

	switch reaction.ReactionType {
	case db.ReactionTypeEmoji:
		// Convert comma-separated emoji into a StringSlice value
		reaction.Reaction = db.StringSlice(strings.Split(strings.TrimSpace(args[3].StringValue()), ","))
		if err := validateEmoji(reaction.Reaction); err != nil {
			return err
		}
	case db.ReactionTypeText:
		// Create a StringSlice with the desired text inside
		reaction.Reaction = db.StringSlice{args[3].StringValue()}
	}

	err := db.AddReaction(i.GuildID, reaction)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, "Successfully added reaction!")
}

// reactionsListCmd handles the `/reactions list` command.
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
		sb.WriteString(reaction.Reaction.String())
		sb.WriteString("\" _(")
		sb.WriteString(string(reaction.ReactionType))
		sb.WriteString(")_\n")
	}

	return util.RespondEphemeral(s, i.Interaction, sb.String())
}

// reactionsDeleteCmd handles the `/reactions delete` command.
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

// reactionsExcludeCmd handles the `/reactions exclude` command.
func reactionsExcludeCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Make sure the user has the manage expressions permission
	// in case a role/member override allows someone else to use it
	if i.Member.Permissions&discordgo.PermissionManageEmojis == 0 {
		return errors.New("you do not have permission to exclude channels")
	}

	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	channel := args[0].ChannelValue(s)

	var match string
	if len(args) == 2 {
		match = args[1].StringValue()
	}

	err := db.ReactionsExclude(i.GuildID, match, channel.ID)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully excluded %s from receiving reactions", channel.Mention()))
}

// reactionsUnexcludeCmd handles the `/reactions unexclude` command.
func reactionsUnexcludeCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Make sure the user has the manage expressions permission
	// in case a role/member override allows someone else to use it
	if i.Member.Permissions&discordgo.PermissionManageEmojis == 0 {
		return errors.New("you do not have permission to unexclude channels")
	}

	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	channel := args[0].ChannelValue(s)

	var match string
	if len(args) == 2 {
		match = args[1].StringValue()
	}

	err := db.ReactionsUnexclude(i.GuildID, match, channel.ID)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully unexcluded %s from receiving reactions", channel.Mention()))
}

// validateEmoji checks if the given slice of emoji is valid.
// If an invalid emoji is found, it returns an error.
func validateEmoji(s db.StringSlice) error {
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
		if _, ok := emoji.Parse(s[i]); !ok {
			return fmt.Errorf("invalid reaction emoji: %s", s[i])
		}
	}
	return nil
}
