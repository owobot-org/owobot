package plugins

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/dop251/goja"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/systems/commands"
	"go.elara.ws/owobot/internal/systems/plugins/builtins"
	"go.elara.ws/owobot/internal/util"
)

func Init(s *discordgo.Session) error {
	if err := loadEnabled(); err != nil {
		return err
	}

	commands.Register(s, pluginCmd, &discordgo.ApplicationCommand{
		Name:        "plugin",
		Description: "Interact with the plugins on this server",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "run",
				Description: "Run a plugin command",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "cmd",
						Description:  "The plugin command to run",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "help",
				Description: "See how to use a plugin command",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "cmd",
						Description:  "The plugin command to help with",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},
	})

	commands.Register(s, pluginadmCmd, &discordgo.ApplicationCommand{
		Name:                     "pluginadm",
		Description:              "Manage dynamic plugins for your server",
		DefaultMemberPermissions: util.Pointer[int64](discordgo.PermissionManageServer),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "List all available plugins",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "enable",
				Description: "Enable a plugin in this guild",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "plugin",
						Description: "The name of the plugin to enable",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "disable",
				Description: "Disable a plugin in this guild",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "plugin",
						Description: "The name of the plugin to disable",
						Required:    true,
					},
				},
			},
		},
	})

	s.AddHandler(handleAutocomplete)
	s.AddHandler(handlePluginEvent)
	return nil
}

// Load recursively loads plugins from the given directory.
func Load(dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".js" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		vm := lockableRuntime{&sync.Mutex{}, goja.New()}
		vm.SetFieldNameMapper(lowerCamelNameMapper{})

		api := &owobotAPI{vm: vm, path: path}

		err = errors.Join(
			vm.GlobalObject().Set("owobot", api),
			vm.GlobalObject().Set("discord", builtins.Constants),
			vm.GlobalObject().Set("print", fmt.Println),
		)
		if err != nil {
			return err
		}

		_, err = vm.RunScript(path, string(data))
		if err != nil {
			return err
		}

		if !api.PluginInfo.IsValid() {
			log.Warn("Plugin info not provided, skipping.").Str("path", path).Send()
			return nil
		}

		prev, _ := db.GetPlugin(api.PluginInfo.Name)

		err = db.AddPlugin(api.PluginInfo)
		if err != nil {
			return err
		}

		err = builtins.Register(vm.Runtime, api.PluginInfo.Name, api.PluginInfo.Version)
		if err != nil {
			return err
		}

		Plugins = append(Plugins, Plugin{
			Info:     api.PluginInfo,
			Commands: api.Commands,
			VM:       vm,
			api:      api,
		})

		if api.Init != nil {
			callableInit, ok := goja.AssertFunction(api.Init)
			if !ok {
				log.Warn("Init value is not callable, ignoring.").Str("plugin", api.PluginInfo.Name).Send()
				return nil
			}

			_, err = callableInit(vm.ToValue(api), vm.ToValue(prev))
			if err != nil {
				return fmt.Errorf("%s init: %w", api.PluginInfo.Name, err)
			}
		}

		return nil
	})
}
