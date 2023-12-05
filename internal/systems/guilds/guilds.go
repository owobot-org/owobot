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

// onGuildCreate adds a guild to the database if it doesn't already exist
func onGuildCreate(s *discordgo.Session, gc *discordgo.GuildCreate) {
	err := db.CreateGuild(gc.ID)
	if err != nil {
		log.Warn("Error creating guild").Err(err).Send()
		return
	}
}

// guildSync makes sure all the guilds the bot is in
// exist in the database. If not, it adds them.
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
