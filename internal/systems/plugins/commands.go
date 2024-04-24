package plugins

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dop251/goja"
	"github.com/kballard/go-shellquote"
	"go.elara.ws/owobot/internal/util"
)

// pluginCmd handles the `/plugin` command and routes it to the correct subcommand.
func pluginCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	switch name := data.Options[0].Name; name {
	case "list":
		return listCmd(s, i)
	case "enable":
		return enableCmd(s, i)
	case "disable":
		return disableCmd(s, i)
	default:
		return fmt.Errorf("unknown plugin subcommand: %s", name)
	}
}

// listCmd handles the `/plugin list` command.
func listCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	sb := strings.Builder{}
	for _, plugin := range Plugins {
		sb.WriteString(plugin.Info.Name)
		sb.WriteString(" (")
		sb.WriteString(plugin.Info.Version)
		sb.WriteString(`): "`)
		sb.WriteString(plugin.Info.Desc)
		sb.WriteByte('"')
		if pluginEnabled(i.GuildID, plugin.Info.Name) {
			sb.WriteString(" *")
		}
		sb.WriteByte('\n')
	}
	return util.RespondEphemeral(s, i.Interaction, sb.String())
}

// enableCmd handles the `/plugin enable` command.
func enableCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	pluginName := data.Options[0].Options[0].StringValue()

	plugin, ok := findPlugin(pluginName)
	if !ok {
		return fmt.Errorf("no such plugin: %q", pluginName)
	}

	err := enablePlugin(i.GuildID, pluginName)
	if err != nil {
		return err
	}

	if plugin.api.OnEnable != nil {
		callable, ok := goja.AssertFunction(plugin.api.OnEnable)
		if !ok {
			return fmt.Errorf("onEnable value is not callable")
		}

		_, err := callable(plugin.VM.ToValue(plugin.api), plugin.VM.ToValue(i.GuildID))
		if err != nil {
			return fmt.Errorf("%s onEnable: %w", plugin.Info.Name, err)
		}
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully enabled the %q plugin!", pluginName))
}

// disableCmd handles the `/plugin disable` command.
func disableCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	pluginName := data.Options[0].Options[0].StringValue()

	plugin, ok := findPlugin(pluginName)
	if !ok {
		return fmt.Errorf("no such plugin: %q", pluginName)
	}

	err := disablePlugin(i.GuildID, pluginName)
	if err != nil {
		return err
	}

	if plugin.api.OnDisable != nil {
		callable, ok := goja.AssertFunction(plugin.api.OnDisable)
		if !ok {
			return fmt.Errorf("onDisable value is not callable")
		}

		plugin.VM.Lock()
		defer plugin.VM.Unlock()

		_, err := callable(plugin.VM.ToValue(plugin.api), plugin.VM.ToValue(i.GuildID))
		if err != nil {
			return fmt.Errorf("%s onDisable: %w", plugin.Info.Name, err)
		}
	}

	return util.RespondEphemeral(s, i.Interaction, fmt.Sprintf("Successfully disabled the %q plugin", pluginName))
}

// phelpCmd handles the `/phelp` command.
func phelpCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	cmdStr := data.Options[0].StringValue()

	args, err := shellquote.Split(cmdStr)
	if err != nil {
		return err
	}

	for _, plugin := range Plugins {
		if !pluginEnabled(i.GuildID, plugin.Info.Name) {
			continue
		}

		cmd, _, ok := findCmd(plugin.Commands, args)
		if !ok {
			continue
		}

		for _, perm := range cmd.Permissions {
			if i.Member.Permissions&perm == 0 {
				return errors.New("you don't have permission to execute this command")
			}
		}

		sb := strings.Builder{}
		sb.WriteString("Usage: `")
		sb.WriteString(cmdStr)
		if usage := cmd.usage(); usage != "" {
			sb.WriteString(" " + usage)
		}
		sb.WriteByte('`')

		sb.WriteString("\n\n")
		sb.WriteString("Description:\n```text\n")
		sb.WriteString(cmd.Desc)
		sb.WriteString("\n```\n")

		if len(cmd.Subcommands) > 0 {
			sb.WriteString("Subcommands:\n")
			for _, subcmd := range cmd.Subcommands {
				sb.WriteString("- `")
				sb.WriteString(subcmd.Name)
				if usage := subcmd.usage(); usage != "" {
					sb.WriteString(" " + usage)
				}
				sb.WriteString("`: `")
				sb.WriteString(subcmd.Desc)
				sb.WriteString("`\n")
			}

		}

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{{
					Title:       "Command `" + cmd.Name + "`",
					Description: sb.String(),
				}},
			},
		})
	}

	return fmt.Errorf("command not found: %q", args[0])
}

// prunCmd handles the `/prunCmd` command.
func prunCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData()
	cmdStr := data.Options[0].StringValue()

	args, err := shellquote.Split(cmdStr)
	if err != nil {
		return err
	}

	for _, plugin := range Plugins {
		if !pluginEnabled(i.GuildID, plugin.Info.Name) {
			continue
		}

		cmd, newArgs, ok := findCmd(plugin.Commands, args)
		if !ok {
			continue
		}

		for _, perm := range cmd.Permissions {
			if i.Member.Permissions&perm == 0 {
				return errors.New("you don't have permission to execute this command")
			}
		}

		callable, ok := goja.AssertFunction(cmd.OnExec)
		if !ok {
			return fmt.Errorf("value in onExec is not callable")
		}

		plugin.VM.Lock()
		_, err = callable(
			plugin.VM.ToValue(cmd),
			plugin.VM.ToValue(s),
			plugin.VM.ToValue(i),
			plugin.VM.ToValue(newArgs),
		)
		plugin.VM.Unlock()
		return err
	}

	return fmt.Errorf("command not found: %q", args[0])
}

func findPlugin(name string) (Plugin, bool) {
	for _, plugin := range Plugins {
		if plugin.Info.Name == name {
			return plugin, true
		}
	}
	return Plugin{}, false
}

func findCmd(cmds []Command, args []string) (Command, []string, bool) {
	if len(args) == 0 {
		return Command{}, nil, false
	}

	for _, cmd := range cmds {
		if args[0] != cmd.Name {
			continue
		}

		if len(cmd.Subcommands) != 0 && len(args) > 1 {
			subcmd, newArgs, ok := findCmd(cmd.Subcommands, args[1:])
			if ok {
				return subcmd, newArgs, true
			}
		}

		return cmd, args[1:], true
	}
	return Command{}, nil, false
}
