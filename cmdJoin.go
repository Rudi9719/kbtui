// +build !rm_basic_commands allcommands joincmd

package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"samhofi.us/x/keybase"
)

func init() {
	command := Command{
		Cmd:         []string{"j", "join"},
		Description: "Join a channel",
		Help:        "",
		Exec:        cmdJoin,
	}

	RegisterCommand(command)
}

func cmdJoin(g *gocui.Gui, cmd []string) {
	stream = false
	if len(cmd) == 3 {
		channel.MembersType = keybase.TEAM
		channel.Name = cmd[1]
		channel.TopicName = cmd[2]
		printToView(g, "Feed", fmt.Sprintf("You are joining: @%s#%s", channel.Name, channel.TopicName))
		clearView(g, "Chat")
		go populateChat(g)
	} else if len(cmd) == 2 {
		channel.MembersType = keybase.USER
		channel.Name = cmd[1]
		channel.TopicName = ""
		printToView(g, "Feed", fmt.Sprintf("You are joining: @%s", channel.Name))
		clearView(g, "Chat")
		go populateChat(g)
	} else {
		printToView(g, "Feed", fmt.Sprintf("To join a team use %sjoin <team> <channel>", cmdPrefix))
		printToView(g, "Feed", fmt.Sprintf("To join a PM use %sjoin <user>", cmdPrefix))
	}
}
