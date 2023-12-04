package members

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
)

var (
	mu        = sync.RWMutex{}
	inviteMap = map[string]*discordgo.Invite{}
)

// populateInviteMap gets invites from all the guilds this bot is in
// and adds them to the invite map.
func populateInviteMap(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		invites, err := s.GuildInvites(guild.ID)
		if err != nil {
			log.Warn("Error getting invites for guild").Str("guild-id", guild.ID).Str("task", "populate-invites").Send()
			continue
		}
		addToMap(invites)
	}
	log.Info("Invite map populated").Send()
}

// addToMap adds a slice of invites to the map
func addToMap(invites []*discordgo.Invite) {
	mu.Lock()
	defer mu.Unlock()
	for _, invite := range invites {
		inviteMap[invite.Code] = invite
	}
}

// addOneToMap adds a single invite to the map
func addOneToMap(invite *discordgo.Invite) {
	mu.Lock()
	defer mu.Unlock()
	inviteMap[invite.Code] = invite
}

// findLastUsedInvites attempts to detect the invites that potentially might've been used last
// in order to figure out what invite a user used to join.
func findLastUsedInvites(s *discordgo.Session, guildID string) ([]string, error) {
	invites, err := s.GuildInvites(guildID)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, invite := range invites {
		mu.RLock()
		oldInvite, ok := inviteMap[invite.Code]
		mu.RUnlock()
		if ok {
			if oldInvite.Uses != invite.Uses {
				out = append(out, invite.Code)
			}
		} else {
			if invite.Uses > 0 {
				out = append(out, invite.Code)
			}
		}
		addOneToMap(invite)
	}
	return out, nil
}
