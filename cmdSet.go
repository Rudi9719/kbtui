// +build !rm_basic_commands allcommands setcmd

package main

import (
	"fmt"
	"strings"
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

func cmdSet(cmd []string) {
	if len(cmd) < 2 {
		printToView("Feed", "No config value specified")
		return
	}
	if len(cmd) < 3 {
		switch cmd[1] {
		case "load":
			printToView("Feed", "Load values from file?")
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
