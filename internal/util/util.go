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

package util

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
)

// Pointer returns a pointer to v. This is useful
// for creating pointers to literal values, as Go doesn't
// allow that by default.
func Pointer[T any](v T) *T {
	return &v
}

// FormatJucheTime formats the given time in Juche calendar format,
// using pyongyang time if it's available, and otherwise UTC.
func FormatJucheTime(t time.Time) string {
	tz, err := time.LoadLocation("Asia/Pyongyang")
	if err != nil {
		tz = time.UTC
	}
	t = t.In(tz)
	return fmt.Sprintf("%02d:%02d %02d-%02d Juche %d", t.Hour(), t.Minute(), t.Day(), t.Month(), t.Year()-1911)
}

// RespondEphemeral responds to an interaction with an ephemeral message.
func RespondEphemeral(s *discordgo.Session, i *discordgo.Interaction, content string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: content,
		},
	})
}

// Respond responds to an interaction with a message.
func Respond(s *discordgo.Session, i *discordgo.Interaction, content string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

// InteractionErrorHandler takes an InteractionCreate event handler that returns an error,
// and returns a regular handler that handles any error by responding with an ephemeral
// message and logging to stderr.
func InteractionErrorHandler(name string, fn func(s *discordgo.Session, i *discordgo.InteractionCreate) error) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := fn(s, i)
		if err != nil {
			log.Warn("Error in interaction handler").Str("name", name).Err(err).Send()
			err = RespondEphemeral(s, i.Interaction, "ERROR: "+err.Error())
			if err != nil {
				log.Warn("Error responding with error").Err(err).Send()
			}
		}
	}
}
