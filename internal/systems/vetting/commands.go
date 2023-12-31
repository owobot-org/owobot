package vetting

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/systems/eventlog"
	"go.elara.ws/owobot/internal/systems/tickets"
	"go.elara.ws/owobot/internal/util"
)

// vettingCmd handles the `/vetting` command and routes it to the correct subcommand.
func vettingCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	switch name := data.Options[0].Name; name {
	case "role":
		return vettingRoleCmd(s, i)
	case "req_channel":
		return vettingReqChannelCmd(s, i)
	case "welcome_channel":
		return welcomeChannelCmd(s, i)
	case "welcome_msg":
		return welcomeMsgCmd(s, i)
	default:
		return fmt.Errorf("unknown vetting subcommand: %s", name)
	}
}

// vettingRoleCmd handles the `/vetting role` command.
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

// vettingReqChannelCmd handles the `/vetting req_channel` command.
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

// welcomeChannelCmd handles the `/vetting welcome_channel` command.
func welcomeChannelCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options
	channel := args[0].ChannelValue(s)

	err := db.SetWelcomeChannel(i.GuildID, channel.ID)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully set %s as the welcome channel!", channel.Mention()))
}

// welcomeMsgCmd handles the `/vetting welcome_msg` command.
func welcomeMsgCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	args := data.Options[0].Options

	err := db.SetWelcomeMsg(i.GuildID, args[0].StringValue())
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, "Successfully set the welcome message!")
}

// approveCmd handles the `/approve` command.
func approveCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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

	err = db.RemoveVettingReq(i.GuildID, i.Message.ID)
	if err != nil {
		return err
	}

	err = eventlog.Log(s, i.GuildID, eventlog.Entry{
		Title:       "New Member Approved!",
		Description: fmt.Sprintf("**User:** %s\n**Role:** %s\n**Approved By:** %s", user.Mention(), role.Mention(), i.Member.User.Mention()),
		Author:      user,
	})
	if err != nil {
		return err
	}

	err = welcomeUser(s, guild, user)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, "Successfully approved "+user.Mention()+" as "+role.Mention()+"!")
}

func welcomeUser(s *discordgo.Session, guild db.Guild, user *discordgo.User) error {
	if guild.WelcomeChanID != "" && guild.WelcomeMsg != "" {
		msg := strings.Replace(guild.WelcomeMsg, "$user", user.Mention(), 1)
		_, err := s.ChannelMessageSend(guild.WelcomeChanID, msg)
		return err
	}
	return nil
}
