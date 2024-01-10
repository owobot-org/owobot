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
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/systems/eventlog"
	"go.elara.ws/owobot/internal/systems/tickets"
	"go.elara.ws/owobot/internal/util"
)

// onMemberJoin adds the vetting role to a user when they join in order to allow them
// to access the vetting questions
func onMemberJoin(s *discordgo.Session, gma *discordgo.GuildMemberAdd) {
	guild, err := db.GuildByID(gma.GuildID)
	if err != nil {
		log.Warn("Error getting guild from database").Str("guild-id", gma.GuildID).Str("task", "vetting-member-join").Err(err).Send()
		return
	}

	if guild.VettingRoleID == "" || guild.VettingReqChanID == "" {
		err = welcomeUser(s, guild, gma.Member.User)
		if err != nil {
			log.Warn("Error welcoming user").Err(err)
		}
		return
	}

	err = s.GuildMemberRoleAdd(gma.GuildID, gma.User.ID, guild.VettingRoleID)
	if err != nil {
		log.Warn("Error assigning vetting role to new user").Str("guild-id", gma.GuildID).Str("task", "vetting-member-join").Err(err).Send()
		return
	}
}

// onMakeVettingMsg deletes and reposts a message with a vetting request button
// that allows users to request vetting.
func onMakeVettingMsg(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	msg := data.Resolved.Messages[data.TargetID]

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg.Content,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Request Vetting",
						Style:    discordgo.SuccessButton,
						Disabled: false,
						Emoji:    &discordgo.ComponentEmoji{Name: clipboardEmoji},
						CustomID: "vetting-req",
					},
				}},
			},
		},
	})
	if err != nil {
		return err
	}

	return s.ChannelMessageDelete(msg.ChannelID, msg.ID)
}

// onVettingRequest handles sends vetting requests in the vetting request channel
func onVettingRequest(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type != discordgo.InteractionMessageComponent {
		return nil
	}

	data := i.MessageComponentData()

	if data.CustomID != "vetting-req" {
		return nil
	}

	_, err := db.VettingReqMsgID(i.GuildID, i.Member.User.ID)
	if err == nil {
		return errors.New("you've already sent a vetting request")
	}

	guild, err := db.GuildByID(i.GuildID)
	if err != nil {
		return err
	}

	if guild.VettingRoleID == "" || guild.VettingReqChanID == "" {
		return nil
	}

	if !slices.Contains(i.Member.Roles, guild.VettingRoleID) {
		return errors.New("you do not have the vetting role")
	}

	embed := &discordgo.MessageEmbed{
		Title: "Vetting Request",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    i.Member.User.Username,
			IconURL: i.Member.User.AvatarURL(""),
		},
		Description: "Accept the vetting request to create a ticket, or reject it to kick the user.",
	}

	eventlog.AddTimeToEmbed(guild.TimeFormat, embed)

	msg, err := s.ChannelMessageSendComplex(guild.VettingReqChanID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Accept",
					Emoji:    &discordgo.ComponentEmoji{Name: checkEmoji},
					Style:    discordgo.SuccessButton,
					CustomID: "vetting-accept:" + i.Member.User.ID,
				},
				discordgo.Button{
					Label:    "Reject",
					Emoji:    &discordgo.ComponentEmoji{Name: crossEmoji},
					Style:    discordgo.DangerButton,
					CustomID: "vetting-reject:" + i.Member.User.ID,
				},
			}},
		},
	})
	if err != nil {
		return err
	}

	err = db.AddVettingReq(i.GuildID, i.Member.User.ID, msg.ID)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, "Successfully sent your vetting request!")
}

// onVettingResponse handles responses to vetting requests. If the user was accepted,
// it creates a vetting ticket for them. If they were rejected, it kicks them from the server.
func onVettingResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type != discordgo.InteractionMessageComponent {
		return nil
	}

	data := i.MessageComponentData()

	resType, userID, ok := strings.Cut(data.CustomID, ":")
	if !ok {
		return nil
	}

	if resType != "vetting-accept" && resType != "vetting-reject" {
		return nil
	}

	guild, err := db.GuildByID(i.GuildID)
	if err != nil {
		return err
	}

	executor := i.Member
	member, err := cache.Member(s, i.GuildID, userID)
	if err != nil {
		return err
	}

	switch resType {
	case "vetting-accept":
		channelID, err := tickets.Open(s, i.GuildID, member.User, executor.User)
		if err != nil {
			return err
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Vetting Request Accepted!",
			Description: fmt.Sprintf("This vetting request has been accepted and a vetting ticket has been created at <#%s>.\n\n**Accepted By:** <@%s>", channelID, executor.User.ID),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    member.User.Username,
				IconURL: member.User.AvatarURL(""),
			},
		}

		eventlog.AddTimeToEmbed(guild.TimeFormat, embed)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			return err
		}
	case "vetting-reject":
		err = s.GuildMemberDeleteWithReason(i.GuildID, member.User.ID, "Vetting request rejected")
		if err != nil {
			return err
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Vetting Request Rejected",
			Description: fmt.Sprintf("This vetting request has been rejected and <@%s> has been kicked from the server.\n\n**Rejected By:** <@%s>", member.User.ID, executor.User.ID),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    member.User.Username,
				IconURL: member.User.AvatarURL(""),
			},
		}

		eventlog.AddTimeToEmbed(guild.TimeFormat, embed)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// onMemberLeave handles users leaving the server. It closes any tickets they might've had open.
func onMemberLeave(s *discordgo.Session, gmr *discordgo.GuildMemberRemove) {
	msgID, err := db.VettingReqMsgID(gmr.GuildID, gmr.Member.User.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return
	} else if err != nil {
		log.Error("Error getting vetting request ID after member leave").Str("user-id", gmr.Member.User.ID).Err(err).Send()
		return
	}

	guild, err := db.GuildByID(gmr.GuildID)
	if err != nil {
		log.Error("Error getting guild").Str("guild-id", gmr.GuildID).Err(err).Send()
		return
	}

	if guild.VettingReqChanID != "" {
		err = s.ChannelMessageDelete(guild.VettingReqChanID, msgID)
		if err != nil {
			log.Error("Error deleting vetting request message after member leave").Str("msg-id", msgID).Err(err).Send()
		}
	}

	err = db.RemoveVettingReq(gmr.GuildID, msgID)
	if err != nil {
		log.Error("Error removing vetting request after member leave").Str("user-id", gmr.Member.User.ID).Err(err).Send()
	}
}
