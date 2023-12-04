package members

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/limiter"
)

// limiters contains the rate limiters for each event type
var limiters = map[string]*limiter.Limiter{
	"channel_delete": limiter.New(8, 10, time.Minute),
	"kick":           limiter.New(8, 10, time.Minute),
	"ban":            limiter.New(5, 7, 5*time.Minute),
}

// handleRatelimit handles rate limiting for a given guild and user ID.
// It decrements the token count for the event, then checks if a warning
// or kick is required, and performs that if needed.
func handleRatelimit(s *discordgo.Session, limit, guildID, userID string) error {
	l, ok := limiters[limit]
	if !ok || userID == s.State.User.ID {
		return nil
	}

	ch, err := s.UserChannelCreate(userID)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", guildID, userID)

	var serverName string
	guild, err := s.State.Guild(guildID)
	if err != nil || guild.Name == "" {
		serverName = "the server"
	} else {
		serverName = guild.Name
	}

	l.Decrement(key)
	if l.IsWarning(key) {
		_, err = s.ChannelMessageSend(ch.ID, fmt.Sprintf("<@%s> WARNING: You are near the `%s` limit. If you reach it, you'll be kicked from %s.", userID, limit, serverName))
		if err != nil {
			return err
		}
	} else if l.IsDepleted(key) {
		// We ignore the error here because even if the message can't be sent,
		// we want to kick the user from the server anyway.
		_, _ = s.ChannelMessageSend(ch.ID, fmt.Sprintf("<@%s> You've reached the `%s` limit and have been kicked from %s.", userID, limit, serverName))

		err = s.GuildMemberDelete(guildID, userID, discordgo.WithAuditLogReason(fmt.Sprintf("Exceeded `%s` rate limit", limit)))
		if err != nil {
			return err
		}
	}

	return nil
}
