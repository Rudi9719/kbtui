// +ignore
// +build type_commands autoreactcmd

package main

import (
	"samhofi.us/x/keybase"
)

func init() {
	command := TypeCommand{
		Cmd:         []string{"text"},
		Name:        "AutoReact",
		Description: "Automatically reacts to every incoming message with an emoji",
		Exec:        tcmdAutoReact,
	}

	RegisterTypeCommand(command)
}

func tcmdAutoReact(m keybase.ChatAPI) {
	msgID := m.Msg.ID
	channel := m.Msg.Channel
	chat := k.NewChat(channel)
	if m.Msg.Sender.Username == "majortrips" {

		chat.React(msgID, ":+1:")
	}
}
