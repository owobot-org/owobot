package plugins

import (
	"reflect"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// HandlerFunc is an event handler function.
type HandlerFunc func(session *discordgo.Session, data any)

// Handler represents a plugin event handler.
type Handler struct {
	PluginName string
	Func       HandlerFunc
}

var (
	handlersMtx = sync.Mutex{}
	handlerMap  = map[string][]Handler{}
)

// handlePluginEvent handles any discord event we receive and
// routes it to the appropriate plugin handler(s).
func handlePluginEvent(s *discordgo.Session, data any) {
	name := reflect.TypeOf(data).Elem().Name()
	handlers, ok := handlerMap[name]
	if !ok {
		return
	}

	for _, h := range handlers {
		if !pluginEnabled(eventGuildID(data), h.PluginName) {
			continue
		}

		h.Func(s, data)
	}
}

// eventGuildID uses reflection to get the guild ID from an event
func eventGuildID(event any) string {
	evt := reflect.ValueOf(event)

	for evt.Kind() == reflect.Pointer {
		evt = evt.Elem()
	}

	if evt.Kind() != reflect.Struct {
		return ""
	}

	if id := evt.FieldByName("GuildID"); id.IsValid() {
		return id.String()
	} else if guild := evt.FieldByName("Guild"); guild.IsValid() {
		if id := guild.FieldByName("ID"); id.IsValid() {
			return id.String()
		}
	}

	return ""
}

// handleAutocomplete handles autocomplete events for the /plugin run command.
func handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	data := i.ApplicationCommandData()
	if data.Name != "prun" && data.Name != "phelp" {
		return
	}

	cmdStr := data.Options[0].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: getAllChoices(i.GuildID, cmdStr, i.Member),
		},
	})
}

// getAllChoices gets possible command strings for each plugin and converts them
// to Discord command options.
func getAllChoices(guildID, partial string, member *discordgo.Member) (out []*discordgo.ApplicationCommandOptionChoice) {
	for _, plugin := range Plugins {
		if !pluginEnabled(guildID, plugin.Info.Name) {
			continue
		}
		out = append(out, getChoiceStrs(partial, "", plugin.Commands, member)...)
	}
	return out
}

// getChoiceStrs recursively looks through every command in cmds,
// and generates a list of strings to use as autocomplete options.
func getChoiceStrs(partial, prefix string, cmds []Command, member *discordgo.Member) []*discordgo.ApplicationCommandOptionChoice {
	if len(cmds) == 0 {
		return nil
	}

	partial = strings.TrimSpace(partial)
	var out []*discordgo.ApplicationCommandOptionChoice

	for _, cmd := range cmds {
		for _, perm := range cmd.Permissions {
			if member.Permissions&perm == 0 {
				continue
			}
		}

		sub := getChoiceStrs(strings.TrimPrefix(partial, cmd.Name), cmd.Name+" ", cmd.Subcommands, member)
		out = append(out, sub...)

		if cmd.OnExec == nil {
			continue
		}

		qualifiedCmd := prefix + cmd.Name

		if strings.Contains(qualifiedCmd, partial) {
			out = append(out, &discordgo.ApplicationCommandOptionChoice{
				Name:  qualifiedCmd,
				Value: qualifiedCmd,
			})
		}
	}

	return out
}
