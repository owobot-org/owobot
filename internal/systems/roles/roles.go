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
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/systems/commands"
	"go.elara.ws/owobot/internal/util"
)

func Init(s *discordgo.Session) error {
	s.AddHandler(util.InteractionErrorHandler("on-role-btn", onRoleButton))

	commands.Register(s, reactionRolesCmd, &discordgo.ApplicationCommand{
		Name:                     "reaction_roles",
		Description:              "Manage reaction roles",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageServer),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "new_category",
				Description: "Create a new reaction role category",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "name",
						Description: "The name of the category",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "description",
						Description: "The description of the reaction role category",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove_category",
				Description: "Remove a reaction role category",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "name",
						Description: "The category which should be removed",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Create a new reaction role",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "category",
						Description: "The category to which the reaction role will be added",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "The role to add",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "emoji",
						Description: "The emoji to use for the role",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a reaction role",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "category",
						Description: "The category from which to remove the reaction role",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "The role to remove",
						Required:    true,
					},
				},
			},
		},
	})

	commands.Register(s, neopronounCmd, &discordgo.ApplicationCommand{
		Name:        "neopronoun",
		Description: "Assign a neopronoun role to yourself",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "neopronoun",
				Description: "The neopronouns to assign to you",
				Required:    true,
			},
		},
	})

	return nil
}

// onRoleButton handles users clicking a role reaction button. It checks if they have
// the role the button is codes for, and if they do, it removes it. Otherwise, it
// assigns it to them.
func onRoleButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type != discordgo.InteractionMessageComponent {
		return nil
	}

	data := i.MessageComponentData()

	buttonID, roleID, ok := strings.Cut(data.CustomID, ":")
	if !ok || buttonID != "role" {
		return nil
	}

	if slices.Contains(i.Member.Roles, roleID) {
		err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleID)
		if err != nil {
			return err
		}
		return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Unassigned role <@&%s>", roleID))
	} else {
		err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID)
		if err != nil {
			return err
		}
		return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully assigned role <@&%s> to you", roleID))
	}
}
