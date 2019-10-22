// +build !rm_basic_commands allcommands reactcmd

package main

import (
	"strconv"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"react", "r", "+"},
		Description: "$messageID $reaction - React to a message (messageID is optional)",
		Help:        "",
		Exec:        cmdReact,
	}

	RegisterCommand(command)
}

func cmdReact(cmd []string) {
	if len(cmd) > 2 {
		reactToMessageID(cmd[1], strings.Join(cmd[2:], " "))
	} else if len(cmd) == 2 {
		reactToMessage(cmd[1])
	}

}

func reactToMessage(reaction string) {
	doReact(lastMessage.ID, reaction)
}
func reactToMessageID(messageID string, reaction string) {
	ID, _ := strconv.Atoi(messageID)
	doReact(ID, reaction)
}
func doReact(messageID int, reaction string) {
	chat := k.NewChat(channel)
	_, err := chat.React(messageID, reaction)
	if err != nil {
		printToView("Feed", "There was an error reacting to the message.")
	}
}
