// +build !rm_basic_commands allcommands editcmd

package main

import (
	"fmt"
	"strconv"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"edit", "e"},
		Description: "$messageId - Edit a message (messageID is optional)",
		Help:        "",
		Exec:        cmdEdit,
	}

	RegisterCommand(command)
}

func cmdEdit(cmd []string) {
	var messageId int
	if len(cmd) == 2 {
		messageId, _ = strconv.Atoi(cmd[1])
		printToView("Input", fmt.Sprintf("/edit %d Type edit here",messageId))
		return
	}
	if len(cmd) < 3 {
		printToView("Feed", "Not enough options for Edit")
		return
	}
	messageId, _ = strconv.Atoi(cmd[1])
	chat := k.NewChat(channel)
	newMessage := strings.Join(cmd[2:], " ")
	_, err := chat.Edit(messageId,newMessage)
	if err != nil {
		printToView("Feed", fmt.Sprintf("Error editing message %d, %+v", messageId, err))
	}
	

}

