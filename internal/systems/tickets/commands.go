package tickets

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/util"
)

// ticketCmd handles the `/ticket` command.
func ticketCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	chID, err := Open(s, i.GuildID, i.Member.User, i.Member.User)
	if err != nil {
		return err
	}
	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully opened a ticket at <#%s>!", chID))
}

// modTicketCmd handles the `/mod_ticket` command.
func modTicketCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	chID, err := Open(s, i.GuildID, data.Options[0].UserValue(s), i.Member.User)
	if err != nil {
		return err
	}
	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully opened a ticket at <#%s>!", chID))
}

// closeTicketCmd handles the `/close_ticket` command.
func closeTicketCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	user := data.Options[0].UserValue(s)
	err := Close(s, i.GuildID, user, i.Member.User)
	if err != nil {
		return err
	}
	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully closed ticket for <@%s>", user.ID))
}

// ticketCategoryCmd handles the `/ticket_category` command.
func ticketCategoryCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	category := data.Options[0].ChannelValue(s)
	err := db.SetTicketCategory(i.GuildID, category.ID)
	if err != nil {
		return err
	}
	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully set the ticket category to `%s`!", category.Name))
}
