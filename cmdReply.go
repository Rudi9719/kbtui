// +build !rm_basic_commands allcommands replycmd

package main

import (
	"strconv"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"reply", "re"},
		Description: "$messageId $response - Reply to a message",
		Help:        "",
		Exec:        cmdReply,
	}

	RegisterCommand(command)
}

func cmdReply(cmd []string) {
	chat := k.NewChat(channel)
	messageId, err := strconv.Atoi(cmd[1])
	_, err = chat.Reply(messageId, strings.Join(cmd[2:], " "))
	if err != nil {
		printToView("Feed", "There was an error with your reply.")
	}
	return
}
