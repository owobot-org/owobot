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

package db

func AddTicket(guildID, userID, channelID string) error {
	_, err := db.Exec("INSERT OR ABORT INTO tickets (guild_id, user_id, channel_id) VALUES (?, ?, ?)", guildID, userID, channelID)
	return err
}

func TicketChannelID(guildID, userID string) (string, error) {
	var out string
	row := db.QueryRowx("SELECT channel_id FROM tickets WHERE user_id = ? AND guild_id = ?", userID, guildID)
	err := row.Scan(&out)
	return out, err
}

func RemoveTicket(guildID, userID string) error {
	_, err := db.Exec("DELETE FROM tickets WHERE user_id = ? AND guild_id = ?", userID, guildID)
	return err
}
