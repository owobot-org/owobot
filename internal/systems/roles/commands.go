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

package roles

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/emoji"
	"go.elara.ws/owobot/internal/util"
)

// reactionRolesCmd handles the `/reaction_roles` command and routes it to the correct subcommand.
func reactionRolesCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()

	switch name := data.Options[0].Name; name {
	case "new_category":
		return reactionRolesNewCategoryCmd(s, i)
	case "remove_category":
		return reactionRolesRemoveCategoryCmd(s, i)
	case "add":
		return reactionRolesAddCmd(s, i)
	case "remove":
		return reactionRolesRemoveCmd(s, i)
	default:
		return fmt.Errorf("unknown reaction_roles subcommand: %s", name)
	}
}

// reactionRolesNewCategoryCmd handles the `/reaction_roles new_category` command.
func reactionRolesNewCategoryCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	rrc := db.ReactionRoleCategory{
		Name: args[0].StringValue(),
	}

	if len(args) > 1 {
		rrc.Description = args[1].StringValue()
	}

	msg, err := s.ChannelMessageSendEmbed(i.ChannelID, &discordgo.MessageEmbed{
		Title:       rrc.Name,
		Description: rrc.Description,
	})
	if err != nil {
		return err
	}

	rrc.MsgID = msg.ID
	err = db.AddReactionRoleCategory(i.ChannelID, rrc)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully added a new reaction role category called `%s`!", rrc.Name))
}

// reactionRolesRemoveCategoryCmd handles the `/reaction_roles remove_category` command.
func reactionRolesRemoveCategoryCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	name := args[0].StringValue()

	rrc, err := db.GetReactionRoleCategory(i.ChannelID, name)
	if err != nil {
		return err
	}

	err = s.ChannelMessageDelete(rrc.ChannelID, rrc.MsgID)
	if err != nil {
		return err
	}

	err = db.DeleteReactionRoleCategory(i.ChannelID, name)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Removed reaction role category `%s`", args[0].StringValue()))
}

// reactionRolesAddCmd handles the `/reaction_roles add` command.
func reactionRolesAddCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	category := args[0].StringValue()
	role := args[1].RoleValue(s, i.GuildID)
	emojiStr := args[2].StringValue()

	_, ok := emoji.Parse(emojiStr)
	if !ok {
		return fmt.Errorf("invalid reaction role emoji: %s", emojiStr)
	}

	err := db.AddReactionRole(i.ChannelID, category, emojiStr, role)
	if err != nil {
		return err
	}

	err = updateReactionRoleCategoryMsg(s, i.ChannelID, category)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Added reaction role %s to `%s`", role.Mention(), category))
}

// reactionRolesRemoveCmd handles the `/reaction_roles remove` command.
func reactionRolesRemoveCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	category := args[0].StringValue()
	role := args[1].RoleValue(s, i.GuildID)

	err := db.DeleteReactionRole(i.ChannelID, category, role)
	if err != nil {
		return err
	}

	err = updateReactionRoleCategoryMsg(s, i.ChannelID, category)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Removed reaction role %s from `%s`", role.Mention(), category))
}

var neopronounValidationRegex = regexp.MustCompile(`^[a-z]+(/[a-z]+)+$`)

// neopronounCmd handles the `/neopronoun` command.
func neopronounCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	name := data.Options[0].StringValue()
	name = strings.ToLower(name)

	if !neopronounValidationRegex.MatchString(name) {
		return fmt.Errorf("invalid neopronoun: `%s`", name)
	}

	roles, err := cache.Roles(s, i.GuildID)
	if err != nil {
		return err
	}

	var roleID string
	for _, role := range roles {
		// Skip this role if it provides any permissions, so that
		// we don't accidentally grant the member any extra permissions
		if role.Permissions != 0 {
			continue
		}

		if role.Name == name {
			roleID = role.ID
			break
		}
	}

	if roleID == "" {
		role, err := s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
			Name:        name,
			Mentionable: util.Pointer(false),
			Permissions: util.Pointer[int64](0),
		})
		if err != nil {
			return err
		}
		roleID = role.ID
	}

	if slices.Contains(i.Member.Roles, roleID) {
		err = s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleID)
		if err != nil {
			return err
		}
		return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Unassigned the `%s` role", name))
	} else {
		err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID)
		if err != nil {
			return err
		}
		return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully assigned the `%s` role to you!", name))
	}
}

// updateReactionRoleCategoryMsg updates a reaction role category message
func updateReactionRoleCategoryMsg(s *discordgo.Session, channelID, category string) error {
	rrc, err := db.GetReactionRoleCategory(channelID, category)
	if err != nil {
		return err
	}

	var sb strings.Builder
	if rrc.Description != "" {
		sb.WriteString(rrc.Description)
		sb.WriteString("\n\n")
	}

	var (
		components []discordgo.MessageComponent
		currentRow discordgo.ActionsRow
	)

	for i, roleID := range rrc.Roles {
		// Action rows can only contain 5 elements,
		// so we create a new row if we reach a multiple
		// of 5.
		if i > 0 && i%5 == 0 {
			components = append(components, currentRow)
			currentRow = discordgo.ActionsRow{}
		}

		e, ok := emoji.Parse(rrc.Emoji[i])
		if !ok {
			return fmt.Errorf("invalid reaction role emoji: %s", rrc.Emoji[i])
		}

		sb.WriteString(rrc.Emoji[i])
		sb.WriteString(" - <@&")
		sb.WriteString(roleID)
		sb.WriteString(">\n")

		currentRow.Components = append(currentRow.Components, discordgo.Button{
			CustomID: "role:" + roleID,
			Style:    discordgo.SecondaryButton,
			Emoji: &discordgo.ComponentEmoji{
				Name: e.Name,
				ID:   e.ID,
			},
		})
	}

	if len(currentRow.Components) > 0 {
		components = append(components, currentRow)
	}

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: channelID,
		ID:      rrc.MsgID,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       rrc.Name,
				Description: sb.String(),
			},
		},
		Components: components,
	})
	return err
}
