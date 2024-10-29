package plugins

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
	"go.elara.ws/logger/log"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/util"
)

// Plugins is a list of plugins
var Plugins []Plugin

// Plugin represents an owobot plugin
type Plugin struct {
	Info     db.PluginInfo
	Commands []Command
	Loop     *eventloop.EventLoop
	api      *owobotAPI
}

// Command represents a plugin command
type Command struct {
	Name        string
	Desc        string
	Usage       goja.Value
	OnExec      goja.Value
	Permissions []int64
	Subcommands []Command
}

func (c Command) usage() string {
	if c.Usage == nil {
		return ""
	} else {
		return c.Usage.String()
	}
}

type owobotAPI struct {
	PluginInfo db.PluginInfo
	Init       goja.Value
	OnEnable   goja.Value
	OnDisable  goja.Value
	Commands   []Command

	path string
	loop *eventloop.EventLoop
}

func (oa *owobotAPI) Enabled(guildID string) bool {
	return pluginEnabled(guildID, oa.PluginInfo.Name)
}

func (oa *owobotAPI) Respond(s *discordgo.Session, i *discordgo.Interaction, content string) error {
	return util.Respond(s, i, content)
}

func (oa *owobotAPI) RespondEphemeral(s *discordgo.Session, i *discordgo.Interaction, content string) error {
	return util.RespondEphemeral(s, i, content)
}

// On adds an event handler function for the given event type
func (oa *owobotAPI) On(eventType string, fn goja.Value) {
	if !oa.PluginInfo.IsValid() {
		log.Warn("No plugin information provided, ignoring handler registration.").Str("path", oa.path).Send()
		return
	}

	callable, ok := goja.AssertFunction(fn)
	if !ok {
		log.Warn("Value passed to handler registrar is not a function, ignoring.").
			Str("plugin", oa.PluginInfo.Name).
			Str("event-type", eventType).
			Send()
		return
	}

	handlersMtx.Lock()
	defer handlersMtx.Unlock()

	oa.loop.RunOnLoop(func(vm *goja.Runtime) {
		this := vm.ToValue(oa)

		handlerMap[eventType] = append(handlerMap[eventType], Handler{
			PluginName: oa.PluginInfo.Name,
			Func: func(s *discordgo.Session, data any) {
				_, err := callable(this, vm.ToValue(s), vm.ToValue(data))
				if err != nil {
					log.Error("Exception thrown in plugin function").
						Str("plugin", oa.PluginInfo.Name).
						Str("event-type", eventType).
						Err(err).
						Send()
				}
			},
		})
	})
}
