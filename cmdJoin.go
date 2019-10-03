// +build !rm_basic_commands allcommands joincmd

package main

import (
	"fmt"

	"samhofi.us/x/keybase"
)

func init() {
	command := Command{
		Cmd:         []string{"join", "j"},
		Description: "Join a channel",
		Help:        "",
		Exec:        cmdJoin,
	}

	RegisterCommand(command)
}

func cmdJoin(cmd []string) {
	stream = false
	if len(cmd) == 3 {
		channel.MembersType = keybase.TEAM
		channel.Name = cmd[1]
		channel.TopicName = cmd[2]
		printToView("Feed", fmt.Sprintf("You are joining: @%s#%s", channel.Name, channel.TopicName))
		clearView("Chat")
		go populateChat()
	} else if len(cmd) == 2 {
		channel.MembersType = keybase.USER
		channel.Name = cmd[1]
		channel.TopicName = ""
		printToView("Feed", fmt.Sprintf("You are joining: @%s", channel.Name))
		clearView("Chat")
		go populateChat()
	} else {
		printToView("Feed", fmt.Sprintf("To join a team use %sjoin <team> <channel>", cmdPrefix))
		printToView("Feed", fmt.Sprintf("To join a PM use %sjoin <user>", cmdPrefix))
	}
}
