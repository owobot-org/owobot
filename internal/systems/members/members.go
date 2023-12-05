package members

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/systems/eventlog"
)

func Init(s *discordgo.Session) error {
	go populateInviteMap(s)
	s.AddHandler(onMemberAdd)
	s.AddHandler(onMemberUpdate)
	s.AddHandler(onMemberLeave)
	s.AddHandler(onChannelDelete)
	return nil
}

// onMemberAdd attempts to detect which invite(s) were used to invite the user
// and logs the member join.
func onMemberAdd(s *discordgo.Session, gma *discordgo.GuildMemberAdd) {
	invites, err := findLastUsedInvites(s, gma.GuildID)
	if err != nil {
		log.Warn("Error finding last used invite").Err(err).Send()
	}

	code := "Unknown"
	if len(invites) > 0 {
		code = strings.Join(invites, " or ")
	}

	err = eventlog.Log(s, gma.GuildID, eventlog.Entry{
		Title:       "New Member Joined!",
		Description: fmt.Sprintf("**User:**\n%s\n**ID:**\n%s\n**Invite Code:**\n%s", gma.Member.User.Mention(), gma.Member.User.ID, code),
		Author:      gma.Member.User,
	})
	if err != nil {
		log.Warn("Error sending member joined log").Str("member", gma.Member.User.Username).Err(err).Send()
	}
}

// onMemberUpdate logs member updates, such as roles being assigned or removed
func onMemberUpdate(s *discordgo.Session, gmu *discordgo.GuildMemberUpdate) {
	if gmu.BeforeUpdate == nil || gmu.Member == nil {
		return
	}

	if !slices.Equal(gmu.BeforeUpdate.Roles, gmu.Member.Roles) {
		var added, removed []string
		for _, newRole := range gmu.Member.Roles {
			if !slices.Contains(gmu.BeforeUpdate.Roles, newRole) {
				added = append(added, fmt.Sprintf("<@&%s>", newRole))
			}
		}
		for _, oldRole := range gmu.BeforeUpdate.Roles {
			if !slices.Contains(gmu.Member.Roles, oldRole) {
				removed = append(removed, fmt.Sprintf("<@&%s>", oldRole))
			}
		}

		err := eventlog.Log(s, gmu.GuildID, eventlog.Entry{
			Title: "Roles Updated",
			Description: fmt.Sprintf(
				"**User:** %s\n**Added:** %s\n**Removed:** %s",
				gmu.Member.User.Mention(),
				strings.Join(added, " "),
				strings.Join(removed, " "),
			),
			Author: gmu.Member.User,
		})
		if err != nil {
			log.Warn("Error roles updated log").Str("member", gmu.Member.User.Username).Err(err).Send()
		}
	}
}

// onMemberLeave logs member leave events and handles bans and kicks
func onMemberLeave(s *discordgo.Session, gmr *discordgo.GuildMemberRemove) {
	err := handleBanOrKick(s, gmr)
	if err != nil {
		log.Warn("Error logging ban or kick").Str("member", gmr.Member.User.Username).Err(err).Send()
	}

	err = eventlog.Log(s, gmr.GuildID, eventlog.Entry{
		Title:       "Member Left",
		Description: fmt.Sprintf("**User:**\n%s\n**ID:**\n%s", gmr.Member.User.Mention(), gmr.Member.User.ID),
		Author:      gmr.Member.User,
	})
	if err != nil {
		log.Warn("Error sending member left log").Str("member", gmr.Member.User.Username).Err(err).Send()
	}
}

// onChannelDelete attempts to detect the user responsible for a channel deletion
// and logs it. It also handles rate limiting for channel delete events.
func onChannelDelete(s *discordgo.Session, cd *discordgo.ChannelDelete) {
	if cd.Type == discordgo.ChannelTypeDM || cd.Type == discordgo.ChannelTypeGroupDM {
		return
	}

	auditLog, err := s.GuildAuditLog(cd.GuildID, "", "", int(discordgo.AuditLogActionChannelDelete), 5)
	if err != nil {
		log.Error("Error getting audit log").Err(err).Send()
		return
	}

	for _, entry := range auditLog.AuditLogEntries {
		// If the deleted channel isn't the one this event is for,
		// skip it.
		if entry.TargetID != cd.ID {
			continue
		}

		// If the bot deleted the channel, we don't care about this event
		if entry.UserID == s.State.User.ID {
			return
		}

		err = handleRatelimit(s, "channel_delete", cd.GuildID, entry.UserID)
		if err != nil {
			log.Error("Error handling rate limit").Err(err).Send()
		}

		member, err := cache.Member(s, cd.GuildID, entry.UserID)
		if err != nil {
			log.Error("Error getting member").Err(err).Send()
			return
		}

		err = eventlog.Log(s, cd.GuildID, eventlog.Entry{
			Title:       "Channel Deleted",
			Description: fmt.Sprintf("**Name:** `%s`\n**Deleted By:** %s", cd.Name, member.User.Mention()),
			Author:      member.User,
		})
		if err != nil {
			log.Warn("Error sending channel deleted log").Str("channel", cd.Name).Err(err).Send()
			return
		}

		return
	}
}

// handleBanOrKick attempts to detect the user responsible for a ban or kick, and
// logs it. It also handles rate limiting for bans and kicks.
func handleBanOrKick(s *discordgo.Session, gmr *discordgo.GuildMemberRemove) error {
	auditLog, err := s.GuildAuditLog(gmr.GuildID, "", "", 0, 5)
	if err != nil {
		return err
	}

	for _, entry := range auditLog.AuditLogEntries {
		// If there's no action type or the user isn't the one this
		// event is for, skip it.
		if entry.ActionType == nil || entry.TargetID != gmr.User.ID {
			continue
		}

		switch *entry.ActionType {
		case discordgo.AuditLogActionMemberBanAdd:
			executor, err := cache.Member(s, gmr.GuildID, entry.UserID)
			if err != nil {
				return err
			}

			err = eventlog.Log(s, gmr.GuildID, eventlog.Entry{
				Title:       "User banned",
				Description: fmt.Sprintf("**Target:** %s\n**Banned by:** %s\n**Reason:** %s", gmr.User.Mention(), executor.User.Mention(), entry.Reason),
				Author:      gmr.User,
			})
			if err != nil {
				return err
			}

			return handleRatelimit(s, "ban", gmr.GuildID, executor.User.ID)
		case discordgo.AuditLogActionMemberKick:
			executor, err := cache.Member(s, gmr.GuildID, entry.UserID)
			if err != nil {
				return err
			}

			err = eventlog.Log(s, gmr.GuildID, eventlog.Entry{
				Title:       "User kicked",
				Description: fmt.Sprintf("**Target:** %s\n**Kicked by:** %s\n**Reason:** %s", gmr.User.Mention(), executor.User.Mention(), entry.Reason),
				Author:      gmr.User,
			})
			if err != nil {
				return err
			}

			return handleRatelimit(s, "kick", gmr.GuildID, executor.User.ID)
		}
	}

	return nil
}
