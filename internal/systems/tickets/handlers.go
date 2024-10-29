package tickets

import (
	"database/sql"
	"errors"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
)

// onMemberLeave closes any tickets a user had open when they leave
func onMemberLeave(s *discordgo.Session, gmr *discordgo.GuildMemberRemove) {
	// If the user had a ticket open when they left, make sure to close it.
	err := Close(s, gmr.GuildID, gmr.User, s.State.User)
	if errors.Is(err, sql.ErrNoRows) {
		// If the error is ErrNoRows, the user didn't have a ticket, so just return
		return
	} else if err != nil {
		log.Warn("Error removing ticket after user left").Err(err).Send()
		return
	}
}
