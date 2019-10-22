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
	chat := k.NewChat(channel)
	if len(cmd) == 2 || len(cmd) == 1 {
		if len(cmd) == 2 {
			messageId, _ = strconv.Atoi(cmd[1])
		} else if lastMessage.ID != 0 {
			if lastMessage.Type != "text" {
				printToView("Feed", "Last message isn't editable (is it an edit?)")
				return
			}
			messageId = lastMessage.ID
		} else {
			printToView("Feed", "No message to edit")
			return
		}
		origMessage, _ := chat.ReadMessage(messageId)
		if origMessage.Result.Messages[0].Msg.Content.Type != "text" {
			printToView("Feed", fmt.Sprintf("%+v", origMessage))
			return
		}
		if origMessage.Result.Messages[0].Msg.Sender.Username != k.Username {
			printToView("Feed", "You cannot edit another user's messages.")
			return
		}
		editString := origMessage.Result.Messages[0].Msg.Content.Text.Body
		clearView("Edit")
		popupView("Edit")
		printToView("Edit", fmt.Sprintf("/e %d %s", messageId, editString))
		setViewTitle("Edit", fmt.Sprintf(" Editing message %d ", messageId))
		return
	}
	if len(cmd) < 3 {
		printToView("Feed", "Not enough options for Edit")
		return
	}
	messageId, _ = strconv.Atoi(cmd[1])
	newMessage := strings.Join(cmd[2:], " ")
	_, err := chat.Edit(messageId, newMessage)
	if err != nil {
		printToView("Feed", fmt.Sprintf("Error editing message %d, %+v", messageId, err))
	}

}
