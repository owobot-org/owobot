package members

import "github.com/bwmarrin/discordgo"

func Init(s *discordgo.Session) error {
	go populateInviteMap(s)
	s.AddHandler(onMemberAdd)
	s.AddHandler(onMemberUpdate)
	s.AddHandler(onMemberLeave)
	s.AddHandler(onChannelDelete)
	return nil
}
