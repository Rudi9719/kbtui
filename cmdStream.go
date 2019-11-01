// +build !rm_basic_commands allcommands streamcmd

package main

import (
	"fmt"
)

func init() {
	command := Command{
		Cmd:         []string{"stream", "s"},
		Description: "- Stream all incoming messages",
		Help:        "",
		Exec:        cmdStream,
	}

	RegisterCommand(command)
}

func cmdStream(cmd []string) {
	stream = true
	channel.Name = ""

	printInfo("You are now viewing the formatted stream")
	setViewTitle("Input", fmt.Sprintf(" Stream - Not in a chat. %sj to join ", config.Basics.CmdPrefix))
	clearView("Chat")
}
