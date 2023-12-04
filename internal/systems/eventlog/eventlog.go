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

package eventlog

import (
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/systems/commands"
	"go.elara.ws/owobot/internal/util"
)

func Init(s *discordgo.Session) error {
	commands.Register(s, eventlogCmd, &discordgo.ApplicationCommand{
		Name:                     "eventlog",
		Description:              "Manage the event log",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageServer),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "channel",
				Description: "Set the event log channel",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:         "channel",
						Description:  "The channel for the event log",
						Type:         discordgo.ApplicationCommandOptionChannel,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						Required:     true,
					},
				},
			},
			{
				Name:        "ticket_channel",
				Description: "Set the ticket log channel",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:         "channel",
						Description:  "The channel for the ticket log",
						Type:         discordgo.ApplicationCommandOptionChannel,
						ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						Required:     true,
					},
				},
			},
		},
	})

	return nil
}

func eventlogCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	switch name := data.Options[0].Name; name {
	case "channel":
		return channelCmd(s, i)
	case "ticket_channel":
		return ticketChannelCmd(s, i)
	default:
		return fmt.Errorf("unknown eventlog subcommand: %s", name)
	}
}

func channelCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Get the subcommand options
	args := i.ApplicationCommandData().Options[0].Options

	c := args[0].ChannelValue(s)
	err := db.SetLogChannel(i.GuildID, c.ID)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully set event log channel to <#%s>!", c.ID))
}

func ticketChannelCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Get the subcommand options
	args := i.ApplicationCommandData().Options[0].Options

	c := args[0].ChannelValue(s)
	err := db.SetTicketLogChannel(i.GuildID, c.ID)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully set ticket log channel to <#%s>!", c.ID))
}

type Entry struct {
	Title       string
	Description string
	ImageURL    string
	Author      *discordgo.User
}

func Log(s *discordgo.Session, guildID string, e Entry) error {
	guild, err := db.GuildByID(guildID)
	if err != nil {
		return err
	}

	if guild.LogChanID == "" {
		return nil
	}

	embed := &discordgo.MessageEmbed{
		Title:       e.Title,
		Description: e.Description,
		Footer: &discordgo.MessageEmbedFooter{
			Text: util.FormatJucheTime(time.Now()),
		},
	}

	if e.Author != nil {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    e.Author.Username,
			IconURL: e.Author.AvatarURL(""),
		}
	}

	if e.ImageURL != "" {
		embed.Image = &discordgo.MessageEmbedImage{URL: e.ImageURL}
	}

	_, err = s.ChannelMessageSendEmbed(guild.LogChanID, embed)
	return err
}

func TicketMsgLog(s *discordgo.Session, guildID string, msgLog io.Reader) error {
	guild, err := db.GuildByID(guildID)
	if err != nil {
		return err
	}

	if guild.TicketLogChanID == "" {
		return nil
	}

	_, err = s.ChannelFileSend(guild.TicketLogChanID, "log.txt", msgLog)
	return err
}
