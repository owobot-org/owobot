package starboard

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/util"
)

// starboardCmd handles the `/starboard` command and routes it to the correct subcommand.
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

// channelCmd handles the `/starboard channel` command.
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

// starsCmd handles the `/starboard stars` command.
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
