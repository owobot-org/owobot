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
	"errors"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ReactionRoleCategory struct {
	MsgID       string      `db:"msg_id"`
	ChannelID   string      `db:"channel_id"`
	Name        string      `db:"name"`
	Description string      `db:"description"`
	Emoji       StringSlice `db:"emoji"`
	Roles       StringSlice `db:"roles"`
}

func AddReactionRoleCategory(channelID string, rrc ReactionRoleCategory) error {
	_, err := db.Exec(
		"INSERT INTO reaction_role_categories VALUES (?, ?, ?, ?, ?, ?)",
		rrc.MsgID,
		channelID,
		rrc.Name,
		rrc.Description,
		rrc.Emoji,
		rrc.Roles,
	)
	return err
}

func GetReactionRoleCategory(channelID, name string) (*ReactionRoleCategory, error) {
	out := &ReactionRoleCategory{}
	err := db.QueryRowx("SELECT * FROM reaction_role_categories WHERE channel_id = ? AND name = ?", channelID, name).StructScan(out)
	return out, err
}

func DeleteReactionRoleCategory(channelID, name string) error {
	_, err := db.Exec("DELETE FROM reaction_role_categories WHERE name = ? AND channel_id = ?", name, channelID)
	return err
}

func AddReactionRole(channelID, category, emojiStr string, role *discordgo.Role) error {
	if strings.Contains(category, "\x1F") || strings.Contains(emojiStr, "\x1F") {
		return errors.New("reaction roles cannot contain unit separator")
	}

	var emoji, roles StringSlice
	err := db.QueryRow("SELECT emoji, roles FROM reaction_role_categories WHERE name = ? AND channel_id = ?", category, channelID).Scan(&emoji, &roles)
	if err != nil {
		return err
	}

	emoji = append(emoji, strings.TrimSpace(emojiStr))
	roles = append(roles, role.ID)

	_, err = db.Exec(
		"UPDATE reaction_role_categories SET emoji = ?, roles = ? WHERE name = ? AND channel_id = ?",
		emoji,
		roles,
		category,
		channelID,
	)
	return err
}

func DeleteReactionRole(channelID, category string, role *discordgo.Role) error {
	var emoji, roles StringSlice
	err := db.QueryRow("SELECT emoji, roles FROM reaction_role_categories WHERE name = ? AND channel_id = ?", category, channelID).Scan(&emoji, &roles)
	if err != nil {
		return err
	}

	if i := slices.Index(roles, role.ID); i == -1 {
		return nil
	} else {
		emoji = append(emoji[:i], emoji[i+1:]...)
		roles = append(roles[:i], roles[i+1:]...)
	}

	_, err = db.Exec(
		"UPDATE reaction_role_categories SET emoji = ?, roles = ? WHERE name = ? AND channel_id = ?",
		emoji,
		roles,
		category,
		channelID,
	)
	return err
}
