package tickets

import (
	"bytes"
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/systems/commands"
	"go.elara.ws/owobot/internal/systems/eventlog"
	"go.elara.ws/owobot/internal/util"
)

const ticketPermissions = discordgo.PermissionSendMessages | discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory

func Init(s *discordgo.Session) error {
	s.AddHandler(onMemberLeave)

	commands.Register(s, ticketCmd, &discordgo.ApplicationCommand{
		Name:        "ticket",
		Description: "Open a ticket to talk to the mods",
	})

	commands.Register(s, ticketCategoryCmd, &discordgo.ApplicationCommand{
		Name:                     "ticket_category",
		Description:              "Set the category in which to create ticket channels",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageServer),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "category",
				Description:  "The category to put ticket channels in",
				Type:         discordgo.ApplicationCommandOptionChannel,
				ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildCategory},
				Required:     true,
			},
		},
	})

	commands.Register(s, modTicketCmd, &discordgo.ApplicationCommand{
		Name:                     "mod_ticket",
		Description:              "Open a ticket for a user to talk to the mods",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageChannels),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The user to open a ticket for",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
		},
	})

	commands.Register(s, closeTicketCmd, &discordgo.ApplicationCommand{
		Name:                     "close_ticket",
		Description:              "Close a user's ticket",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageChannels),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The user whose ticket to close",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
		},
	})

	return nil
}

// Open opens a new ticket. It checks if a ticket already exists, and if not, creates a new channel for it,
// allows the user it's for to see and send messages in it, adds it to the database, and logs the ticket open.
func Open(s *discordgo.Session, guildID string, user, executor *discordgo.User) (string, error) {
	channelID, err := db.TicketChannelID(guildID, user.ID)
	if err == nil {
		return "", fmt.Errorf("ticket already exists for %s at <#%s>", user.Mention(), channelID)
	}

	if executor == nil {
		executor = s.State.User
	}

	guild, err := db.GuildByID(guildID)
	if err != nil {
		return "", err
	}

	overwrites := []*discordgo.PermissionOverwrite{{
		ID:    user.ID,
		Type:  discordgo.PermissionOverwriteTypeMember,
		Allow: ticketPermissions,
	}}

	if guild.TicketCategoryID != "" {
		category, err := cache.Channel(s, guildID, guild.TicketCategoryID)
		if err != nil {
			log.Error("Error getting ticket category").Err(err).Send()
			// If we can't get the ticket category, set it to empty string
			// so that ChannelCreate doesn't try to use it.
			guild.TicketCategoryID = ""
		} else {
			overwrites = append(overwrites, category.PermissionOverwrites...)
		}
	}

	c, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:                 "ticket-" + user.Username,
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             guild.TicketCategoryID,
		PermissionOverwrites: overwrites,
	})
	if err != nil {
		return "", err
	}

	err = db.AddTicket(guildID, user.ID, c.ID)
	if err != nil {
		return "", err
	}

	return c.ID, eventlog.Log(s, guildID, eventlog.Entry{
		Title:       "New ticket opened!",
		Description: "**Executed by:** " + executor.Mention(),
		Author:      user,
	})
}

// Close closes the given user's ticket. It gets the channel ID of the ticket, logs all the messages
// inside it, deletes the channel, removes the ticket from the database, and logs the ticket close.
func Close(s *discordgo.Session, guildID string, user, executor *discordgo.User) error {
	channelID, err := db.TicketChannelID(guildID, user.ID)
	if err != nil {
		return err
	}

	if executor == nil {
		executor = s.State.User
	}

	guild, err := db.GuildByID(guildID)
	if err != nil {
		return err
	}

	if guild.TicketLogChanID != "" {
		buf, err := getChannelMessageLog(s, channelID)
		if err != nil {
			return err
		}

		if buf != nil {
			err = eventlog.TicketMsgLog(s, guildID, buf)
			if err != nil {
				return err
			}
		}
	}

	_, err = s.ChannelDelete(channelID)
	if err != nil {
		return err
	}

	err = db.RemoveTicket(guildID, user.ID)
	if err != nil {
		return err
	}

	return eventlog.Log(s, guildID, eventlog.Entry{
		Title:       "Ticket Closed",
		Description: "**Executed by:** " + executor.Mention(),
		Author:      user,
	})
}

// getChannelMessageLog generates a log for the given channel. It retrieves all the messages
// inside it and writes them to a buffer.
func getChannelMessageLog(s *discordgo.Session, channelID string) (*bytes.Buffer, error) {
	out := &bytes.Buffer{}

	msgs, err := s.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		return nil, err
	}

	if len(msgs) == 0 {
		return nil, nil
	}

	err = writeMsgs(msgs, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// writeMsgs writes a slice of messages to w.
func writeMsgs(msgs []*discordgo.Message, w io.Writer) error {
	if len(msgs) == 0 {
		return nil
	}

	for i := len(msgs) - 1; i >= 0; i-- {
		_, err := io.WriteString(w, fmt.Sprintf("%s - %s\n", msgs[i].Author.Username, msgs[i].Content))
		if err != nil {
			return err
		}
	}
	return nil
}
