package about

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/systems/commands"
)

const aboutTmpl = `**Copyright © %d owobot contributors**

This program comes with **ABSOLUTELY NO WARRANTY**. This is free software, and you are welcome to redistribute it under certain conditions. See [here](https://www.gnu.org/licenses/agpl-3.0.html) for details.

**Running Commit:**
%s

**Source Code:**
https://gitea.elara.ws/owobot/owobot
**GitHub Mirror:**
https://github.com/owobot-org/owobot`

func Init(s *discordgo.Session) error {
	commands.Register(s, aboutCmd, &discordgo.ApplicationCommand{
		Name:        "about",
		Description: "Information about owobot",
	})
	return nil
}

func aboutCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{{
				Title: "About owobot",
				Description: fmt.Sprintf(
					aboutTmpl,
					time.Now().Year(),
					getCommit(),
				),
			}},
		},
	})
}

func getCommit() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "`<unknown>`"
	}

	commit := "`<unknown>`"
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			commit = "`" + setting.Value + "`"
		case "vcs.modified":
			if setting.Value == "true" {
				commit += " (modified)"
			}
		}
	}

	return commit
}
