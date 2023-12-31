package starboard

import (
	"fmt"
	"mime"
	"net/url"
	"path"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/systems/eventlog"
	"mvdan.cc/xurls/v2"
)

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
				mra.GuildID,
				msg.ChannelID,
				msg.ID,
			),
			Color: embedColor,
		}

		eventlog.AddTimeToEmbed(guild.TimeFormat, embed)

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
	if xurl := xurls.Strict().FindString(msg.Content); xurl != "" {
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
