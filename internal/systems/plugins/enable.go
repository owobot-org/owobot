package plugins

import (
	"fmt"
	"slices"

	"go.elara.ws/owobot/internal/db"
)

var enabled = map[string][]string{}

func loadEnabled() error {
	guilds, err := db.AllGuilds()
	if err != nil {
		return err
	}
	for _, guild := range guilds {
		enabled[guild.ID] = []string(guild.EnabledPlugins)
	}
	return nil
}

func enablePlugin(guildID, pluginName string) error {
	if slices.Contains(enabled[guildID], pluginName) {
		return fmt.Errorf("plugin %q is already enabled", pluginName)
	}
	enabled[guildID] = append(enabled[guildID], pluginName)
	return db.EnablePlugin(guildID, pluginName)
}

func disablePlugin(guildID, pluginName string) error {
	if i := slices.Index(enabled[guildID], pluginName); i > -1 {
		enabled[guildID] = append(enabled[guildID][:i], enabled[guildID][i+1:]...)
	} else {
		return fmt.Errorf("plugin %q is already disabled", pluginName)
	}
	return db.DisablePlugin(guildID, pluginName)
}

func pluginEnabled(guildID, pluginName string) bool {
	if guildID == "" {
		return false
	}
	return slices.Contains(enabled[guildID], pluginName)
}
