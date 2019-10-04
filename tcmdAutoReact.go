// +build type_commands autoreactcmd

package main

import (
	"samhofi.us/x/keybase"
)

func init() {
	command := TypeCommand{
		Cmd:         []string{"text"},
		Description: "Automatically reacts to every incoming message with an emoji",
		Exec:        tcmdAutoReact,
	}

	RegisterTypeCommand(command)
}

func tcmdAutoReact(m keybase.ChatAPI) {
	msgID := m.Msg.ID
	channel := m.Msg.Channel
	chat := k.NewChat(channel)
	chat.React(msgID, ":sunglasses:")
}
