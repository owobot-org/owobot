/*
 * owobot - Your server's guardian and entertainer
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

package cache

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// Member gets a discord member from the cache. If it doesn't exist in the cache, it
// gets it from discord and adds it to the cache.
func Member(s *discordgo.Session, guildID, userID string) (*discordgo.Member, error) {
	member, err := s.State.Member(guildID, userID)
	if errors.Is(err, discordgo.ErrStateNotFound) {
		// If the member wasn't found in the state struct,
		// get the member from discord and add it.
		member, err = s.GuildMember(guildID, userID)
		if err != nil {
			return nil, err
		}
		return member, s.State.MemberAdd(member)
	} else if err != nil {
		return nil, err
	}
	return member, nil
}

// Role gets a discord role from the cache. If it doesn't exist in the cache, it
// gets it from discord and adds it to the cache.
func Role(s *discordgo.Session, guildID, roleID string) (*discordgo.Role, error) {
	role, err := s.State.Role(guildID, roleID)
	if errors.Is(err, discordgo.ErrStateNotFound) {
		// If the role wasn't found in the state struct,
		// get the guild roles from discord and add them.
		roles, err := s.GuildRoles(guildID)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			err = s.State.RoleAdd(guildID, role)
			if err != nil {
				return nil, err
			}
		}
		return s.State.Role(guildID, roleID)
	} else if err != nil {
		return nil, err
	}
	return role, nil
}

// Roles gets a list of roles in a discord guild from the cache. If it doesn't
// exist in the cache, it gets it from discord and adds it to the cache.
func Roles(s *discordgo.Session, guildID string) ([]*discordgo.Role, error) {
	guild, err := s.State.Guild(guildID)
	if errors.Is(err, discordgo.ErrStateNotFound) {
		return s.GuildRoles(guildID)
	} else if err != nil {
		return nil, err
	}

	if len(guild.Roles) == 0 {
		roles, err := s.GuildRoles(guildID)
		if err != nil {
			return nil, err
		}
		guild.Roles = roles
		return roles, nil
	}

	return guild.Roles, nil
}
