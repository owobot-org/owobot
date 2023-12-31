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

package starboard

import (
	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/systems/commands"
	"go.elara.ws/owobot/internal/util"
)

const (
	starEmoji  = "\u2b50"
	embedColor = 0xFF5833
)

func Init(s *discordgo.Session) error {
	s.AddHandler(onReaction)

	commands.Register(s, starboardCmd, &discordgo.ApplicationCommand{
		Name:                     "starboard",
		Description:              "Modify starboard settings",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageServer),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "channel",
				Description: "Set the starboard channel",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:         "channel",
						Description:  "The channel to use for the starboard",
						Type:         discordgo.ApplicationCommandOptionChannel,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						Required:     true,
					},
				},
			},
			{
				Name:        "stars",
				Description: "Set the amount of stars for the starboard",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "stars",
						Description: "The amount of stars to require",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Required:    true,
					},
				},
			},
		},
	})

	return nil
}
