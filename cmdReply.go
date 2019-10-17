// +build !rm_basic_commands allcommands replycmd

package main

import (
	"fmt"
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
	if len(cmd) < 2 {
		printToView("Feed", fmt.Sprintf("%sreply $ID - Reply to message $ID", cmdPrefix))
		return
	}
	messageID, err := strconv.Atoi(cmd[1])
	if err != nil {
		printToView("Feed", fmt.Sprintf("There was an error determining message ID %s", cmd[1]))
		return
	}
	_, err = chat.Reply(messageID, strings.Join(cmd[2:], " "))
	if err != nil {
		printToView("Feed", "There was an error with your reply.")
		return
	}
	return
}
