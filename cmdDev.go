// +build !rm_basic_commands allcommands devcmd

package main

import (
	"fmt"
)

func init() {
	command := Command{
		Cmd:         []string{"dev"},
		Description: "- Switch to dev channels",
		Help:        "",
		Exec:        cmdDev,
	}

	RegisterCommand(command)
}

func cmdDev(cmd []string) {
	dev = !dev

	printInfo(fmt.Sprintf("You have toggled the dev flag to %+v", dev))
	clearView("Chat")
}
