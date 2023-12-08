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
)

type Poll struct {
	MsgID        string
	OwnerID      string
	Title        string
	Finished     bool
	OptionEmojis []string
	OptionText   []string
}

func CreatePoll(msgID, ownerID, title string) error {
	_, err := db.Exec("INSERT INTO polls(msg_id, owner_id, title) VALUES (?, ?, ?)", msgID, ownerID, title)
	return err
}

func GetPoll(msgID string) (*Poll, error) {
	var title, ownerID, emojis, text string
	var finished bool
	err := db.QueryRow("SELECT title, owner_id, finished, opt_emojis, opt_text FROM polls WHERE msg_id = ?", msgID).Scan(&title, &ownerID, &finished, &emojis, &text)
	if err != nil {
		return nil, err
	}
	return &Poll{
		MsgID:        msgID,
		OwnerID:      ownerID,
		Title:        title,
		Finished:     finished,
		OptionEmojis: splitOptions(emojis),
		OptionText:   splitOptions(text),
	}, nil
}

func AddPollOptionText(msgID string, text string) error {
	if strings.Contains(text, "\x1F") {
		return errors.New("option string cannot contain unit separator")
	}
	
	var optText string
	err := db.QueryRow("SELECT opt_text FROM polls WHERE msg_id = ?", msgID).Scan(&optText)
	if err != nil {
		return err
	}

	splitText := splitOptions(optText)
	splitText = append(splitText, text)

	_, err = db.Exec("UPDATE polls SET opt_text = ? WHERE msg_id = ?", strings.Join(splitText, "\x1F"), msgID)
	return err
}

func AddPollOptionEmoji(msgID string, emoji string) error {
	if strings.Contains(emoji, "\x1F") {
		return errors.New("emoji string cannot contain unit separator")
	}

	var optEmojis string
	err := db.QueryRow("SELECT opt_emojis FROM polls WHERE msg_id = ?", msgID).Scan(&optEmojis)
	if err != nil {
		return err
	}

	splitEmojis := splitOptions(optEmojis)
	if slices.Contains(splitEmojis, emoji) {
		return errors.New("emojis can only be used once")
	}
	splitEmojis = append(splitEmojis, emoji)

	_, err = db.Exec("UPDATE polls SET opt_emojis = ? WHERE msg_id = ?", strings.Join(splitEmojis, "\x1F"), msgID)
	return err
}

func FinishPoll(msgID string) error {
	_, err := db.Exec("UPDATE polls SET finished = true WHERE msg_id = ?", msgID)
	return err
}

type Vote struct {
	PollMsgID string `db:"poll_msg_id"`
	UserToken string `db:"user_token"`
	Option    int    `db:"option"`
}

func UserVote(msgID, userToken string) (Vote, error) {
	var out Vote
	row := db.QueryRowx("SELECT * FROM votes WHERE poll_msg_id = ? AND user_token = ?", msgID, userToken)
	err := row.StructScan(&out)
	return out, err
}

func AddVote(v Vote) error {
	_, err := db.NamedExec("INSERT OR REPLACE INTO votes (poll_msg_id, user_token, option) VALUES (:poll_msg_id, :user_token, :option)", v)
	return err
}

func VoteAmount(msgID string, option int) (int64, error) {
	var out int64
	err := db.QueryRow("SELECT COUNT(1) FROM votes WHERE poll_msg_id = ? AND option = ?", msgID, option).Scan(&out)
	return out, err
}

func splitOptions(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\x1F")
}
