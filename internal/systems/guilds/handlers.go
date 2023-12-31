package guilds

import (
	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/db"
)

// onGuildCreate listens for when the bot joins a new guild and adds it
// to the database if it doesn't already exist.
func onGuildCreate(s *discordgo.Session, gc *discordgo.GuildCreate) {
	err := db.CreateGuild(gc.ID)
	if err != nil {
		log.Warn("Error creating guild").Err(err).Send()
		return
	}
}
