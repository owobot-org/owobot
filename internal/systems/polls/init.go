package polls

import (
	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/systems/commands"
	"go.elara.ws/owobot/internal/util"
)

func Init(s *discordgo.Session) error {
	s.AddHandler(util.InteractionErrorHandler("poll-add-opt", onPollAddOpt))
	s.AddHandler(util.InteractionErrorHandler("poll-opt-submit", onAddOptModalSubmit))
	s.AddHandler(util.InteractionErrorHandler("poll-finish", onPollFinish))
	s.AddHandler(onPollReaction)
	s.AddHandler(onVote)

	commands.Register(s, pollCmd, &discordgo.ApplicationCommand{
		Name:        "poll",
		Description: "Create a new poll",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "title",
				Description: "The title of the poll",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	})

	return nil
}
