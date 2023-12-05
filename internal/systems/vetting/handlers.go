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
	"time"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/systems/eventlog"
	"go.elara.ws/owobot/internal/systems/tickets"
	"go.elara.ws/owobot/internal/util"
)

// vettingCmd runs the correct subcommand handler for the vetting command
func vettingCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	switch name := data.Options[0].Name; name {
	case "role":
		return vettingRoleCmd(s, i)
	case "req_channel":
		return vettingReqChannelCmd(s, i)
	default:
		return fmt.Errorf("unknown vetting subcommand: %s", name)
	}
}

// vettingRoleCmd sets the vetting role for a guild
func vettingRoleCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options
	role := args[0].RoleValue(s, i.GuildID)

	err := db.SetVettingRoleID(i.GuildID, role.ID)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully set %s as the vetting role!", role.Mention()))
}

// vettingReqChannelCmd sets the vettign request channel for a guild
func vettingReqChannelCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options
	channel := args[0].ChannelValue(s)

	err := db.SetVettingReqChannel(i.GuildID, channel.ID)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully set %s as the vetting request channel!", channel.Mention()))
}

// onMemberJoin adds the vetting role to a user when they join in order to allow them
// to access the vetting questions
func onMemberJoin(s *discordgo.Session, gma *discordgo.GuildMemberAdd) {
	guild, err := db.GuildByID(gma.GuildID)
	if err != nil {
		log.Warn("Error getting guild from database").Str("guild-id", gma.GuildID).Str("task", "vetting-member-join").Err(err).Send()
		return
	}

	if guild.VettingRoleID == "" {
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
						Emoji:    discordgo.ComponentEmoji{Name: clipboardEmoji},
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

	_, err = s.ChannelMessageSendComplex(guild.VettingReqChanID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Vetting Request",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    i.Member.User.Username,
				IconURL: i.Member.User.AvatarURL(""),
			},
			Description: "Accept the vetting request to create a ticket, or reject it to kick the user.",
			Footer: &discordgo.MessageEmbedFooter{
				Text: util.FormatJucheTime(time.Now()),
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Accept",
					Emoji:    discordgo.ComponentEmoji{Name: checkEmoji},
					Style:    discordgo.SuccessButton,
					CustomID: "vetting-accept:" + i.Member.User.ID,
				},
				discordgo.Button{
					Label:    "Reject",
					Emoji:    discordgo.ComponentEmoji{Name: crossEmoji},
					Style:    discordgo.DangerButton,
					CustomID: "vetting-reject:" + i.Member.User.ID,
				},
			}},
		},
	})
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, "Successfully sent your vetting request!")
}

// onApprove approves a user in vetting. It removes their vetting role, assigns a
// role of the approver's choosing, closes the user's vetting ticket, and logs
// the approval.
func onApprove(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guild, err := db.GuildByID(i.GuildID)
	if err != nil {
		return err
	}

	if guild.VettingRoleID == "" {
		return errors.New("vetting role id is not set for this guild")
	}

	data := i.ApplicationCommandData()
	user := data.Options[0].UserValue(s)
	role := data.Options[1].RoleValue(s, i.GuildID)

	_, err = db.TicketChannelID(i.GuildID, user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s has no open ticket", user.Mention())
	}

	roleSetAllowed := false
	for _, roleID := range i.Member.Roles {
		executorRole, err := cache.Role(s, i.GuildID, roleID)
		if err != nil {
			return err
		}
		if executorRole.Position >= role.Position {
			roleSetAllowed = true
			break
		}
	}

	if !roleSetAllowed {
		return errors.New("you don't have permission to approve a user as a role higher than your own")
	}

	err = s.GuildMemberRoleAdd(i.GuildID, user.ID, role.ID)
	if err != nil {
		return err
	}

	err = s.GuildMemberRoleRemove(i.GuildID, user.ID, guild.VettingRoleID)
	if err != nil {
		return err
	}

	err = tickets.Close(s, i.GuildID, user, i.Member.User)
	if err != nil {
		return err
	}

	err = eventlog.Log(s, i.GuildID, eventlog.Entry{
		Title:       "New Member Approved!",
		Description: fmt.Sprintf("User: %s\nRole: %s\nApproved By: %s", user.Mention(), role.Mention(), i.Member.User.Mention()),
		Author:      user,
	})
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, "Successfully approved "+user.Mention()+" as "+role.Mention()+"!")
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

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Vetting Request Accepted!",
						Description: fmt.Sprintf("This vetting request has been accepted and a vetting ticket has been created at <#%s>.\n\n**Accepted By:** <@%s>", channelID, executor.User.ID),
						Author: &discordgo.MessageEmbedAuthor{
							Name:    member.User.Username,
							IconURL: member.User.AvatarURL(""),
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: util.FormatJucheTime(time.Now()),
						},
					},
				},
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

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Vetting Request Rejected",
						Description: fmt.Sprintf("This vetting request has been rejected and <@%s> has been kicked from the server.\n\n**Rejected By:** <@%s>", member.User.ID, executor.User.ID),
						Author: &discordgo.MessageEmbedAuthor{
							Name:    member.User.Username,
							IconURL: member.User.AvatarURL(""),
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: util.FormatJucheTime(time.Now()),
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
