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
	MsgID       string   `db:"msg_id"`
	ChannelID   string   `db:"channel_id"`
	Name        string   `db:"name"`
	Description string   `db:"description"`
	Emoji       []string `db:"emoji"`
	Roles       []string `db:"roles"`
}

func AddReactionRoleCategory(channelID string, rrc ReactionRoleCategory) error {
	_, err := db.Exec(
		"INSERT INTO reaction_role_categories VALUES (?, ?, ?, ?, ?, ?)",
		rrc.MsgID,
		channelID,
		rrc.Name,
		rrc.Description,
		strings.Join(rrc.Emoji, "\x1F"),
		strings.Join(rrc.Roles, "\x1F"),
	)
	return err
}

func GetReactionRoleCategory(channelID, name string) (*ReactionRoleCategory, error) {
	var msgID, description, emoji, roles string
	err := db.QueryRow(
		"SELECT msg_id, description, emoji, roles FROM reaction_role_categories WHERE channel_id = ? AND name = ?",
		channelID,
		name,
	).Scan(&msgID, &description, &emoji, &roles)
	if err != nil {
		return nil, err
	}

	return &ReactionRoleCategory{
		MsgID:       msgID,
		ChannelID:   channelID,
		Name:        name,
		Description: description,
		Emoji:       splitOptions(emoji),
		Roles:       splitOptions(roles),
	}, nil
}

func DeleteReactionRoleCategory(channelID, name string) error {
	_, err := db.Exec("DELETE FROM reaction_role_categories WHERE name = ? AND channel_id = ?", name, channelID)
	return err
}

func AddReactionRole(channelID, category, emoji string, role *discordgo.Role) error {
	if strings.Contains(category, "\x1F") || strings.Contains(emoji, "\x1F") {
		return errors.New("reaction roles cannot contain unit separator")
	}

	var oldEmoji, oldRoles string
	err := db.QueryRow("SELECT emoji, roles FROM reaction_role_categories WHERE name = ? AND channel_id = ?", category, channelID).Scan(&oldEmoji, &oldRoles)
	if err != nil {
		return err
	}

	splitEmoji, splitRoles := splitOptions(oldEmoji), splitOptions(oldRoles)
	splitEmoji = append(splitEmoji, strings.TrimSpace(emoji))
	splitRoles = append(splitRoles, role.ID)

	_, err = db.Exec(
		"UPDATE reaction_role_categories SET emoji = ?, roles = ? WHERE name = ? AND channel_id = ?",
		strings.Join(splitEmoji, "\x1F"),
		strings.Join(splitRoles, "\x1F"),
		category,
		channelID,
	)
	return err
}

func DeleteReactionRole(channelID, category string, role *discordgo.Role) error {
	var oldEmoji, oldRoles string
	err := db.QueryRow("SELECT emoji, roles FROM reaction_role_categories WHERE name = ? AND channel_id = ?", category, channelID).Scan(&oldEmoji, &oldRoles)
	if err != nil {
		return err
	}

	splitEmoji, splitRoles := splitOptions(oldEmoji), splitOptions(oldRoles)
	if i := slices.Index(splitRoles, role.ID); i == -1 {
		return nil
	} else {
		splitEmoji = append(splitEmoji[:i], splitEmoji[i+1:]...)
		splitRoles = append(splitRoles[:i], splitRoles[i+1:]...)
	}

	_, err = db.Exec(
		"UPDATE reaction_role_categories SET emoji = ?, roles = ? WHERE name = ? AND channel_id = ?",
		strings.Join(splitEmoji, "\x1F"),
		strings.Join(splitRoles, "\x1F"),
		category,
		channelID,
	)
	return err
}
