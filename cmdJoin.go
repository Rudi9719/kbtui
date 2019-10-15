// +build !rm_basic_commands allcommands joincmd

package main

import (
	"fmt"
	"samhofi.us/x/keybase"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"join", "j"},
		Description: "$team/user $channel - Join a chat, $user or $team $channel",
		Help:        "",
		Exec:        cmdJoin,
	}

	RegisterCommand(command)
}

func cmdJoin(cmd []string) {
	stream = false
	switch l := len(cmd); l {
	case 3:
		fallthrough
	case 2:
		// if people write it in one singular line, with a `#`
		firstArgSplit := strings.Split(cmd[1], "#")
		channel.Name = strings.Replace(firstArgSplit[0], "@", "", 1)
		joinedName := fmt.Sprintf("@%s", channel.Name)
		if l == 3 || len(firstArgSplit) == 2 {
			channel.MembersType = keybase.TEAM
			if l == 3 {
				channel.TopicName = strings.Replace(cmd[2], "#", "", 1)
			} else {
				channel.TopicName = firstArgSplit[1]
			}
			joinedName = fmt.Sprintf("%s#%s", joinedName, channel.TopicName)
		} else {
			channel.TopicName = ""
			channel.MembersType = keybase.USER
		}
		printToView("Feed", fmt.Sprintf("You are joining: %s", joinedName))
		clearView("Chat")
		viewTitle("Input", fmt.Sprintf(" %s ", joinedName))
		go populateChat()
	default:
		printToView("Feed", fmt.Sprintf("To join a team use %sjoin <team> <channel>", cmdPrefix))
		printToView("Feed", fmt.Sprintf("To join a PM use %sjoin <user>", cmdPrefix))
	}
}
