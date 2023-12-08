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

package vetting

import (
	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/systems/commands"
	"go.elara.ws/owobot/internal/util"
)

const (
	clipboardEmoji = "\U0001f4cb"
	checkEmoji     = "\u2705"
	crossEmoji     = "\u2694\ufe0f"
)

func Init(s *discordgo.Session) error {
	s.AddHandler(onMemberJoin)
	s.AddHandler(util.InteractionErrorHandler("on-vetting-req", onVettingRequest))
	s.AddHandler(util.InteractionErrorHandler("on-vetting-resp", onVettingResponse))
	s.AddHandler(onMemberLeave)

	commands.Register(s, onMakeVettingMsg, &discordgo.ApplicationCommand{
		Name:                     "Make Vetting Message",
		Type:                     discordgo.MessageApplicationCommand,
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageServer),
	})

	commands.Register(s, vettingCmd, &discordgo.ApplicationCommand{
		Name:                     "vetting",
		Description:              "Manage vetting",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageServer),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "role",
				Description: "Set the vetting role",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "role",
						Description: "The role to use for vetting",
						Type:        discordgo.ApplicationCommandOptionRole,
						Required:    true,
					},
				},
			},
			{
				Name:        "req_channel",
				Description: "Set the vetting request channel",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:         "channel",
						Description:  "The channel to use for vetting requests",
						Type:         discordgo.ApplicationCommandOptionChannel,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						Required:     true,
					},
				},
			},
			{
				Name:        "welcome_channel",
				Description: "Set the welcome channel",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:         "channel",
						Description:  "The channel to use for welcoming new users",
						Type:         discordgo.ApplicationCommandOptionChannel,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						Required:     true,
					},
				},
			},
			{
				Name:        "welcome_msg",
				Description: "Set the welcome message",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "msg",
						Description: "The message to welcome new users with",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			},
		},
	})

	commands.Register(s, approveCmd, &discordgo.ApplicationCommand{
		Name:                     "approve",
		Description:              "Approve a member in vetting",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionKickMembers),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "member",
				Description: "The member to approve",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			{
				Name:        "role",
				Description: "The role to approve the member as",
				Type:        discordgo.ApplicationCommandOptionRole,
				Required:    true,
			},
		},
	})

	return nil
}
