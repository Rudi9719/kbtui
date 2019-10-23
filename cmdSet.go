// +build !rm_basic_commands allcommands setcmd

package main

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml"
)

func init() {
	command := Command{
		Cmd:         []string{"set", "config"},
		Description: "Change various settings",
		Help:        "",
		Exec:        cmdSet,
	}

	RegisterCommand(command)
}
func printSetting(cmd []string) {
	switch cmd[1] {
	case "load":
		loadFromToml()
	case "downloadPath":
		printToView("Feed", fmt.Sprintf("Setting for %s -> %s", cmd[1], downloadPath))
	case "outputFormat":
		printToView("Feed", fmt.Sprintf("Setting for %s -> %s", cmd[1], outputFormat))
	case "dateFormat":
		printToView("Feed", fmt.Sprintf("Setting for %s -> %s", cmd[1], dateFormat))
	case "timeFormat":
		printToView("Feed", fmt.Sprintf("Setting for %s -> %s", cmd[1], timeFormat))
	case "cmdPrefix":
		printToView("Feed", fmt.Sprintf("Setting for %s -> %s", cmd[1], cmdPrefix))
	default:
		printToView("Feed", fmt.Sprintf("Unknown config value %s", cmd[1]))
	}

	return
}
func cmdSet(cmd []string) {
	if len(cmd) < 2 {
		printToView("Feed", "No config value specified")
		return
	}
	if len(cmd) < 3 {
		printSetting(cmd)
		return
	}
	switch cmd[1] {
	case "downloadPath":
		if len(cmd) != 3 {
			printToView("Feed", "Invalid download path.")
		}
		downloadPath = cmd[2]
	case "outputFormat":
		outputFormat = strings.Join(cmd[1:], " ")
	case "dateFormat":
		dateFormat = strings.Join(cmd[1:], " ")
	case "timeFormat":
		timeFormat = strings.Join(cmd[1:], " ")
	case "cmdPrefix":
		cmdPrefix = cmd[2]
	default:
		printToView("Feed", fmt.Sprintf("Unknown config value %s", cmd[1]))
	}

}
func loadFromToml() {
	printToView("Feed", fmt.Sprintf("Loading config from toml"))
	config, err := toml.LoadFile("kbtui.tml")
	if err != nil {
		printToView("Feed", fmt.Sprintf("Could not read config file: %+v", err))
		return
	}
	if config.Has("Basics.colorless") {
		colorless = config.Get("Basics.colorless").(bool)
	}
	if config.Has("Basics.downloadPath") {
		downloadPath = config.Get("Basics.downloadPath").(string)
	}
	if config.Has("Basics.cmdPrefix") {
		cmdPrefix = config.Get("Basics.cmdPrefix").(string)
	}
	if config.Has("Formatting.outputFormat") {
		outputFormat = config.Get("Formatting.outputFormat").(string)
	}
	if config.Has("Formatting.dateFormat") {
		dateFormat = config.Get("Formatting.dateFormat").(string)
	}
	if config.Has("Formatting.timeFormat") {
		timeFormat = config.Get("Formatting.timeFormat").(string)
	}
}
