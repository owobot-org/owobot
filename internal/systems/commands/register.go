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

package commands

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/util"
)

var (
	mu   = sync.Mutex{}
	cmds = map[string]CmdFunc{}
	acs  = []*discordgo.ApplicationCommand{}
)

type CmdFunc func(s *discordgo.Session, i *discordgo.InteractionCreate) error

func Init(s *discordgo.Session) error {
	s.AddHandler(onCmd)
	_, err := s.ApplicationCommandBulkOverwrite(s.State.Application.ID, "", acs)
	acs = nil // Allow the ACs to be GC'd
	return err
}

func onCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	mu.Lock()
	defer mu.Unlock()
	cmdFn, ok := cmds[data.Name]
	if !ok {
		return
	}

	err := cmdFn(s, i)
	if err != nil {
		log.Warn("Error in command function").Str("cmd", data.Name).Err(err).Send()
		sendError(s, i.Interaction, err)
		return
	}
}

func Register(s *discordgo.Session, fn CmdFunc, ac *discordgo.ApplicationCommand) {
	// If the DM permission hasn't been explicitly set, assume false
	if ac.DMPermission == nil {
		ac.DMPermission = util.Pointer(false)
	}

	mu.Lock()
	// Skip commands that already exist
	if _, ok := cmds[ac.Name]; ok {
		mu.Unlock()
		return
	}
	cmds[ac.Name] = fn
	mu.Unlock()

	acs = append(acs, ac)
}

func sendError(s *discordgo.Session, i *discordgo.Interaction, serr error) {
	err := util.RespondEphemeral(s, i, "ERROR: "+serr.Error())
	if err != nil {
		log.Warn("Error while trying to send error").Err(err).Send()
		return
	}
}

// commandSync checks if any registered commands have been removed and, if so,
// deletes them.
func commandSync(s *discordgo.Session) error {
	appCmds, err := s.ApplicationCommands(s.State.Application.ID, "")
	if err != nil {
		return err
	}

	deleted := 0
	for _, appCmd := range appCmds {
		if _, ok := cmds[appCmd.Name]; !ok {
			err = s.ApplicationCommandDelete(appCmd.ApplicationID, "", appCmd.ID)
			if err != nil {
				return err
			}
			deleted++
		}
	}

	log.Info("Command sync completed").Int("deleted", deleted).Send()
	return nil
}
