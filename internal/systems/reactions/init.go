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
	"github.com/bwmarrin/discordgo"
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
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "exclude",
				Description: "Exclude a channel from having reactions",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "channel",
						Description: "The channel which shouldn't receive reactions",
						Type:        discordgo.ApplicationCommandOptionChannel,
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildText,
							discordgo.ChannelTypeGuildForum,
						},
						Required: true,
					},
					{
						Name:        "match",
						Description: "The match value to exclude",
						Type:        discordgo.ApplicationCommandOptionString,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "unexclude",
				Description: "Unexclude a channel from having reactions",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "channel",
						Description: "The channel which should receive reactions",
						Type:        discordgo.ApplicationCommandOptionChannel,
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildText,
							discordgo.ChannelTypeGuildForum,
						},
						Required: true,
					},
					{
						Name:        "match",
						Description: "The match value to unexclude",
						Type:        discordgo.ApplicationCommandOptionString,
					},
				},
			},
		},
	})

	return nil
}
