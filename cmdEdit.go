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
		Description: "$messageID - Edit a message (messageID is optional)",
		Help:        "",
		Exec:        cmdEdit,
	}

	RegisterCommand(command)
}

func cmdEdit(cmd []string) {
	var messageID int
	chat := k.NewChat(channel)
	if len(cmd) == 2 || len(cmd) == 1 {
		if len(cmd) == 2 {
			messageID, _ = strconv.Atoi(cmd[1])
		} else if lastMessage.ID != 0 {
			message, _ := chat.ReadMessage(lastMessage.ID)
			lastMessage.Type = message.Result.Messages[0].Msg.Content.Type
			if lastMessage.Type != "text" {
				printError("Last message isn't editable (is it an edit?)")
				return
			}
			messageID = lastMessage.ID
		} else {
			printError("No message to edit")
			return
		}
		origMessage, _ := chat.ReadMessage(messageID)
		if origMessage.Result.Messages[0].Msg.Content.Type != "text" {
			printInfo(fmt.Sprintf("%+v", origMessage))
			return
		}
		if origMessage.Result.Messages[0].Msg.Sender.Username != k.Username {
			printError("You cannot edit another user's messages.")
			return
		}
		editString := origMessage.Result.Messages[0].Msg.Content.Text.Body
		clearView("Edit")
		popupView("Edit")
		printToView("Edit", fmt.Sprintf("/e %d %s", messageID, editString))
		setViewTitle("Edit", fmt.Sprintf(" Editing message %d ", messageID))
		moveCursorToEnd("Edit")
		return
	}
	if len(cmd) < 3 {
		printError("Not enough options for Edit")
		return
	}
	messageID, _ = strconv.Atoi(cmd[1])
	newMessage := strings.Join(cmd[2:], " ")
	_, err := chat.Edit(messageID, newMessage)
	if err != nil {
		printError(fmt.Sprintf("Error editing message %d, %+v", messageID, err))
	}

}
