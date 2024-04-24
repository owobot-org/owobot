package builtins

import (
	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/systems/eventlog"
	"go.elara.ws/owobot/internal/systems/tickets"
)

type eventLogAPI struct{}

func (eventLogAPI) Log(s *discordgo.Session, guildID string, e eventlog.Entry) error {
	return eventlog.Log(s, guildID, e)
}

type ticketsAPI struct{}

func (ticketsAPI) Open(s *discordgo.Session, guildID string, user, executor *discordgo.User) (string, error) {
	return tickets.Open(s, guildID, user, executor)
}

func (ticketsAPI) Close(s *discordgo.Session, guildID string, user, executor *discordgo.User) {
	tickets.Close(s, guildID, user, executor)
}

type cacheAPI struct{}

func (cacheAPI) Channel(s *discordgo.Session, guildID, channelID string) (*discordgo.Channel, error) {
	return cache.Channel(s, guildID, channelID)
}

func (cacheAPI) Member(s *discordgo.Session, guildID, userID string) (*discordgo.Member, error) {
	return cache.Member(s, guildID, userID)
}

func (cacheAPI) Role(s *discordgo.Session, guildID, roleID string) (*discordgo.Role, error) {
	return cache.Role(s, guildID, roleID)
}

func (cacheAPI) Roles(s *discordgo.Session, guildID string) ([]*discordgo.Role, error) {
	return cache.Roles(s, guildID)
}
