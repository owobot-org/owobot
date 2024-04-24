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

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"go.elara.ws/logger"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/systems/about"
	"go.elara.ws/owobot/internal/systems/commands"
	"go.elara.ws/owobot/internal/systems/eventlog"
	"go.elara.ws/owobot/internal/systems/guilds"
	"go.elara.ws/owobot/internal/systems/members"
	"go.elara.ws/owobot/internal/systems/plugins"
	"go.elara.ws/owobot/internal/systems/polls"
	"go.elara.ws/owobot/internal/systems/reactions"
	"go.elara.ws/owobot/internal/systems/roles"
	"go.elara.ws/owobot/internal/systems/starboard"
	"go.elara.ws/owobot/internal/systems/tickets"
	"go.elara.ws/owobot/internal/systems/vetting"
)

func init() {
	log.Logger = logger.NewPretty(os.Stderr)
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Error loading configuration").Err(err).Send()
	}

	err = db.Init(ctx, cfg.DBPath+"?_pragma=busy_timeout(30000)")
	if err != nil {
		log.Fatal("Error initializing database").Err(err).Send()
	}

	s, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatal("Error creating new session").Err(err).Send()
	}

	s.StateEnabled = true
	s.State.TrackMembers = true
	s.State.TrackRoles = true
	s.State.TrackChannels = true
	s.Identify.Intents |= discordgo.IntentMessageContent | discordgo.IntentGuildMembers

	err = s.Open()
	if err != nil {
		log.Fatal("Error opening a connection to discord").Err(err).Send()
	}

	if cfg.Activity.Type != -1 && cfg.Activity.Name != "" {
		err = s.UpdateStatusComplex(discordgo.UpdateStatusData{Activities: []*discordgo.Activity{
			{Type: cfg.Activity.Type, Name: cfg.Activity.Name},
		}})
		if err != nil {
			log.Error("Error updating status").Err(err).Send()
		}
	}

	err = plugins.Load(cfg.PluginDir)
	if err != nil {
		log.Error("Error running plugin file").Err(err).Send()
	}

	initSystems(
		s,
		starboard.Init,
		members.Init,
		guilds.Init,
		tickets.Init,
		eventlog.Init,
		polls.Init,
		vetting.Init,
		reactions.Init,
		roles.Init,
		about.Init,
		plugins.Init,
		commands.Init, // The commands system should always go last
	)

	log.Info("Everything is initialized, the bot is ready!").Send()

	select {
	case <-ctx.Done():
		log.Info("Context canceled, shutting down...").Send()
		s.Close()
		db.Close()
	}
}

func initSystems(s *discordgo.Session, fns ...func(*discordgo.Session) error) {
	for i, fn := range fns {
		err := fn(s)
		if err != nil {
			log.Warn("Error initializing system").Int("index", i).Err(err).Send()
		}
	}
}
