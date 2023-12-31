package commands

import (
	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/util"
)

// onCmd handles any command interaction and routes it to the correct command
// if it was registered using the [Register] function.
func onCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	mu.Lock()
	cmdFn, ok := cmds[data.Name]
	if !ok {
		mu.Unlock()
		return
	}
	mu.Unlock()

	err := cmdFn(s, i)
	if err != nil {
		log.Warn("Error in command function").Str("cmd", data.Name).Err(err).Send()
		sendError(s, i.Interaction, err)
		return
	}
}

// sendError responds to an interaction with an ephemeral message containing an error
func sendError(s *discordgo.Session, i *discordgo.Interaction, serr error) {
	err := util.RespondEphemeral(s, i, "ERROR: "+serr.Error())
	if err != nil {
		log.Warn("Error while trying to send error").Err(err).Send()
		return
	}
}
