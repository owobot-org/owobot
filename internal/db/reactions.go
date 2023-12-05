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

package db

type MatchType string

const (
	MatchTypeRegex    MatchType = "regex"
	MatchTypeContains MatchType = "contains"
)

type ReactionType string

const (
	ReactionTypeEmoji ReactionType = "emoji"
	ReactionTypeText  ReactionType = "text"
)

type Reaction struct {
	GuildID      string       `db:"guild_id"`
	MatchType    MatchType    `db:"match_type"`
	Match        string       `db:"match"`
	ReactionType ReactionType `db:"reaction_type"`
	Reaction     string       `db:"reaction"`
	Chance       int          `db:"chance"`
}

func AddReaction(guildID string, r Reaction) error {
	r.GuildID = guildID
	_, err := db.NamedExec("INSERT INTO reactions VALUES (:guild_id, :match_type, :match, :reaction_type, :reaction, :chance)", r)
	return err
}

func DeleteReaction(guildID string, match string) error {
	_, err := db.Exec("DELETE FROM reactions WHERE guild_id = ? AND match = ?", guildID, match)
	return err
}

func Reactions(guildID string) (rs []Reaction, err error) {
	err = db.Select(&rs, "SELECT * FROM reactions WHERE guild_id = ?", guildID)
	return rs, err
}
