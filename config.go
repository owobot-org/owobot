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
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v10"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Token     string   `env:"TOKEN" toml:"token"`
	DBPath    string   `env:"DB_PATH" toml:"db_path"`
	PluginDir string   `env:"PLUGIN_DIR" toml:"plugin_dir"`
	Activity  Activity `envPrefix:"ACTIVITY_" toml:"activity"`
}

type Activity struct {
	Type discordgo.ActivityType `env:"TYPE" toml:"type"`
	Name string                 `env:"NAME" toml:"name"`
}

func loadConfig() (*Config, error) {
	// Create a new config struct with default values
	cfg := &Config{
		Token:     "",
		DBPath:    "owobot.db",
		PluginDir: "plugins",
		Activity: Activity{
			Type: -1,
			Name: "",
		},
	}

	configPath := os.Getenv("OWOBOT_CONFIG_PATH")
	if configPath == "" {
		configPath = "/etc/owobot/config.toml"
	}

	fl, err := os.Open(configPath)
	if err == nil {
		err = toml.NewDecoder(fl).Decode(cfg)
		if err != nil {
			return nil, err
		}
		fl.Close()
	}

	return cfg, env.ParseWithOptions(cfg, env.Options{Prefix: "OWOBOT_"})
}
