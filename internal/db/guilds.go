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

import (
	"database/sql"
	"errors"
)

type Guild struct {
	ID               string `db:"id"`
	StarboardChanID  string `db:"starboard_chan_id"`
	StarboardStars   int    `db:"starboard_stars"`
	LogChanID        string `db:"log_chan_id"`
	TicketLogChanID  string `db:"ticket_log_chan_id"`
	TicketCategoryID string `db:"ticket_category_id"`
	VettingReqChanID string `db:"vetting_req_chan_id"`
	VettingRoleID    string `db:"vetting_role_id"`
	TimeFormat       string `db:"time_format"`
	WelcomeChanID    string `db:"welcome_chan_id"`
	WelcomeMsg       string `db:"welcome_msg"`
}

func AllGuilds() ([]Guild, error) {
	var out []Guild
	err := db.Select(&out, "SELECT * FROM guilds")
	return out, err
}

func GuildByID(id string) (Guild, error) {
	var out Guild
	err := db.QueryRowx("SELECT * FROM guilds WHERE id = ? LIMIT 1", id).StructScan(&out)
	return out, err
}

func CreateGuild(guildID string) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO guilds (id) VALUES (?)`, guildID)
	return err
}

func SetStarboardChannel(guildID, channelID string) error {
	_, err := db.Exec("UPDATE guilds SET starboard_chan_id = ? WHERE id = ?", channelID, guildID)
	return err
}

func SetStarboardStars(guildID string, stars int64) error {
	_, err := db.Exec("UPDATE guilds SET starboard_stars = ? WHERE id = ?", stars, guildID)
	return err
}

func SetLogChannel(guildID, channelID string) error {
	_, err := db.Exec("UPDATE guilds SET log_chan_id = ? WHERE id = ?", channelID, guildID)
	return err
}

func SetTicketLogChannel(guildID, channelID string) error {
	_, err := db.Exec("UPDATE guilds SET ticket_log_chan_id = ? WHERE id = ?", channelID, guildID)
	return err
}

func SetTicketCategory(guildID, categoryID string) error {
	_, err := db.Exec("UPDATE guilds SET ticket_category_id = ? WHERE id = ?", categoryID, guildID)
	return err
}

func SetVettingReqChannel(guildID, channelID string) error {
	_, err := db.Exec("UPDATE guilds SET vetting_req_chan_id = ? WHERE id = ?", channelID, guildID)
	return err
}

func SetVettingRoleID(guildID, roleID string) error {
	_, err := db.Exec("UPDATE guilds SET vetting_role_id = ? WHERE id = ?", roleID, guildID)
	return err
}

func SetTimeFormat(guildID, timeFmt string) error {
	_, err := db.Exec("UPDATE guilds SET time_format = ? WHERE id = ?", timeFmt, guildID)
	return err
}

func SetWelcomeChannel(guildID, channelID string) error {
	_, err := db.Exec("UPDATE guilds SET welcome_chan_id = ? WHERE id = ?", channelID, guildID)
	return err
}

func SetWelcomeMsg(guildID, msg string) error {
	_, err := db.Exec("UPDATE guilds SET welcome_msg = ? WHERE id = ?", msg, guildID)
	return err
}

func IsVettingMsg(msgID string) (bool, error) {
	var out bool
	err := db.QueryRow("SELECT 1 FROM guild WHERE vetting_msg_id = ?", msgID).Scan(&out)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
