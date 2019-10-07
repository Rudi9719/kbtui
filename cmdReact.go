// +build !rm_basic_commands allcommands reactcmd

package main

import (
	"strconv"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"react", "r", "+"},
		Description: "$messageId $reaction - React to a message (messageID is optional)",
		Help:        "",
		Exec:        cmdReact,
	}

	RegisterCommand(command)
}

func cmdReact(cmd []string) {
	if len(cmd) > 2 {
		reactToMessageId(cmd[1], strings.Join(cmd[2:], " "))
	} else if len(cmd) == 2 {
		reactToMessage(cmd[1])
	}

}

func reactToMessage(reaction string) {
	doReact(lastMessage.ID, reaction)
}
func reactToMessageId(messageId string, reaction string) {
	ID, _ := strconv.Atoi(messageId)
	doReact(ID, reaction)
}
func doReact(messageId int, reaction string) {
	chat := k.NewChat(channel)
	_, err := chat.React(messageId, reaction)
	if err != nil {
		printToView("Feed", "There was an error reacting to the message.")
	}
}
