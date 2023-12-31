package eventlog

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/util"
)

// eventlogCmd handles the `/eventlog` command and routes it to the correct subcommand.
func eventlogCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	switch name := data.Options[0].Name; name {
	case "channel":
		return channelCmd(s, i)
	case "ticket_channel":
		return ticketChannelCmd(s, i)
	case "time_format":
		return timeFormatCmd(s, i)
	default:
		return fmt.Errorf("unknown eventlog subcommand: %s", name)
	}
}

// channelCmd handles the `/eventlog channel` command.
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

// ticketChannelCmd handles the `/eventlog ticket_channel` command.
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

// timeFormatCmd handles the `/eventlog time_format` command
func timeFormatCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Get the subcommand options
	args := i.ApplicationCommandData().Options[0].Options
	timeFmt := args[0].StringValue()

	err := db.SetTimeFormat(i.GuildID, timeFmt)
	if err != nil {
		return err
	}

	return util.RespondEphemeral(s, i.Interaction, "Successfully set the time format!")
}
