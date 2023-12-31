package roles

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/util"
)

// onRoleButton handles users clicking a role reaction button. It checks if they have
// the role the button is codes for, and if they do, it removes it. Otherwise, it
// assigns it to them.
func onRoleButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Type != discordgo.InteractionMessageComponent {
		return nil
	}

	data := i.MessageComponentData()

	buttonID, roleID, ok := strings.Cut(data.CustomID, ":")
	if !ok || buttonID != "role" {
		return nil
	}

	if slices.Contains(i.Member.Roles, roleID) {
		err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleID)
		if err != nil {
			return err
		}
		return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Unassigned role <@&%s>", roleID))
	} else {
		err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID)
		if err != nil {
			return err
		}
		return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully assigned role <@&%s> to you", roleID))
	}
}