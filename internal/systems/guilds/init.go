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

package guilds

import (
	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/db"
)

func Init(s *discordgo.Session) error {
	s.AddHandler(onGuildCreate)
	return guildSync(s)
}

// guildSync looks through all the guilds that the bot is in,
// and if any of them don't exist in the database, it adds them.
func guildSync(s *discordgo.Session) error {
	for _, guild := range s.State.Guilds {
		err := db.CreateGuild(guild.ID)
		if err != nil {
			return err
		}
	}
	log.Info("Guild sync completed").Send()
	return nil
}
