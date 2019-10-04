// +build !rm_basic_commands allcommands editcmd

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
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
		} else {
			messageId = lastMessage.ID
		}

		origMessage, _ := chat.ReadMessage(messageId)
		if origMessage.Result.Messages[0].Msg.Content.Type != "text" {
			printToView("Feed", fmt.Sprintf("%+v", origMessage))
			return
		}
		editString := origMessage.Result.Messages[0].Msg.Content.Text.Body
		g.Update(func(g *gocui.Gui) error {
			inputView, err := g.View("Input")
			if err != nil {
				return err
			} else {
				editString = fmt.Sprintf("/edit %d %s", messageId, editString)
				fmt.Fprintf(inputView, editString)
				viewX, viewY := inputView.Size()
				if len(editString) < viewX {
					viewX = len(editString)
					viewY = 0
				} else {
					viewX = viewX / len(editString)
				}
				inputView.MoveCursor(viewX, viewY, true)
			}
			return nil
		})
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
