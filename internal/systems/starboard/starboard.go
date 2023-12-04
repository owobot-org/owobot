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

package starboard

import (
	"errors"
	"fmt"
	"mime"
	"net/url"
	"path"
	"strings"
	"time"

	"mvdan.cc/xurls"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/db"
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

// starboardCmd calls the correct subcommand handler for the starboard command
func starboardCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	switch name := data.Options[0].Name; name {
	case "channel":
		return channelCmd(s, i)
	case "stars":
		return starsCmd(s, i)
	default:
		return fmt.Errorf("unknown subcommand: %s", name)
	}
}

// channelCmd sets the starboard channel for the guild
func channelCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Get the subcommand options
	args := i.ApplicationCommandData().Options[0].Options

	c := args[0].ChannelValue(s)
	err := db.SetStarboardChannel(i.GuildID, c.ID)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully set starboard channel to <#%s>!", c.ID))
}

// starsCmd sets the amount of stars that trigger the starboard for the guild
func starsCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Get the subcommand options
	args := i.ApplicationCommandData().Options[0].Options

	stars := args[0].IntValue()
	if stars <= 0 {
		return errors.New("star amount must be greater than 0")
	}

	err := db.SetStarboardStars(i.GuildID, stars)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully set the amount of stars required to get on the starboard to %d!", stars))
}

// onReaction detects star reactions, and if the message qualifies for starboard
// based on the guild's settings, it replies to it and adds it to the starboard.
func onReaction(s *discordgo.Session, mra *discordgo.MessageReactionAdd) {
	if mra.Emoji.Name != starEmoji {
		return
	}

	msgExists, err := db.ExistsInStarboard(mra.MessageID)
	if err != nil {
		log.Warn("Error checking if the message exists in the starboard").Err(err).Send()
		return
	}

	// If the message has already been added to the starboard,
	// we can skip it.
	if msgExists {
		return
	}

	guild, err := db.GuildByID(mra.GuildID)
	if err != nil {
		log.Warn("Error getting guild from the database").Str("id", mra.GuildID).Err(err).Send()
		return
	}

	// If the guild has no starboard channel ID set, we can
	// skip this message.
	if guild.StarboardChanID == "" {
		return
	}

	reactions, err := s.MessageReactions(mra.ChannelID, mra.MessageID, starEmoji, guild.StarboardStars, "", "")
	if err != nil {
		log.Warn("Error getting message reactions").Err(err).Send()
		return
	}

	if len(reactions) >= guild.StarboardStars {
		msg, err := s.ChannelMessage(mra.ChannelID, mra.MessageID)
		if err != nil {
			log.Warn("Error getting channel message").Err(err).Send()
			return
		}

		ch, err := s.Channel(mra.ChannelID)
		if err != nil {
			log.Warn("Error getting channel").Err(err).Send()
			return
		}

		_, err = s.ChannelMessageSendReply(
			msg.ChannelID,
			fmt.Sprintf("Congrats %s! You've made it to <#%s>!!", msg.Author.Mention(), guild.StarboardChanID),
			msg.Reference(),
		)
		if err != nil {
			log.Warn("Error sending message reply").Err(err).Send()
			return
		}

		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s - #%s - %s has made it!", starEmoji, ch.Name, msg.Author.Username),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    msg.Author.Username,
				IconURL: msg.Author.AvatarURL(""),
			},
			Description: fmt.Sprintf(
				"[**Jump to Message**](https://discord.com/channels/%s/%s/%s)",
				msg.GuildID,
				msg.ChannelID,
				msg.ID,
			),
			Color: embedColor,
			Footer: &discordgo.MessageEmbedFooter{
				Text: util.FormatJucheTime(time.Now()),
			},
		}

		if imageURL := getImageURL(msg); imageURL != "" {
			// If the message has an image, add it to the embed
			embed.Image = &discordgo.MessageEmbedImage{URL: imageURL}
		}

		if msg.Content != "" {
			// If the message has content, we add it above the
			// jump to message link currently in the embed description.
			embed.Description = fmt.Sprintf(
				"**Message Content**\n%s\n\n%s",
				msg.Content,
				embed.Description,
			)
		}

		_, err = s.ChannelMessageSendEmbed(guild.StarboardChanID, embed)
		if err != nil {
			log.Warn("Error sending starboard message").Err(err).Send()
			return
		}

		err = db.AddToStarboard(mra.MessageID)
		if err != nil {
			log.Warn("Error adding message to starboard").Err(err).Send()
			return
		}
	}
}

// getImageURL looks through the message content and attachments
// to try to find images. If it finds one, it returns the URL.
// Otherwise, it returns an empty string.
func getImageURL(msg *discordgo.Message) string {
	if xurl := xurls.Strict.FindString(msg.Content); xurl != "" {
		u, err := url.Parse(xurl)
		if err == nil {
			mt := mime.TypeByExtension(path.Ext(u.Path))
			if strings.HasPrefix(mt, "image/") {
				return xurl
			}
		}
	}

	for _, attachment := range msg.Attachments {
		if strings.HasPrefix(attachment.ContentType, "image/") {
			return attachment.URL
		}
	}

	return ""
}
