/*
 * owobot - The coolest Discord bot ever written
 * Copyright (C) 2023 owobot Contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package roles

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/owobot/internal/cache"
	"go.elara.ws/owobot/internal/util"
)

var neopronounValidationRegex = regexp.MustCompile(`^[a-z]+(/[a-z]+)+$`)

// neopronounCmd assigns a neopronoun role to the user that ran it.
func neopronounCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	name := data.Options[0].StringValue()
	name = strings.ToLower(name)

	if !neopronounValidationRegex.MatchString(name) {
		return fmt.Errorf("invalid neopronoun: `%s`", name)
	}

	roles, err := cache.Roles(s, i.GuildID)
	if err != nil {
		return err
	}

	var roleID string
	for _, role := range roles {
		// Skip this role if it provides any permissions, so that
		// we don't accidentally grant the member any extra permissions
		if role.Permissions != 0 {
			continue
		}

		if role.Name == name {
			roleID = role.ID
			break
		}
	}

	if roleID == "" {
		role, err := s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
			Name:        name,
			Mentionable: util.Pointer(false),
			Permissions: util.Pointer[int64](0),
		})
		if err != nil {
			return err
		}
		roleID = role.ID
	}

	if slices.Contains(i.Member.Roles, roleID) {
		err = s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleID)
		if err != nil {
			return err
		}
		return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Unassigned the `%s` role", name))
	} else {
		err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID)
		if err != nil {
			return err
		}
		return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully assigned the `%s` role to you!", name))
	}
}
