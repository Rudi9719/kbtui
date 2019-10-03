// +build !rm_basic_commands allcommands reactcmd

package main

import "strconv"

func init() {
	command := Command{
		Cmd:         []string{"react", "r", "+"},
		Description: "React to a message",
		Help:        "",
		Exec:        cmdReact,
	}

	RegisterCommand(command)
}

func cmdReact(cmd []string) {
	if len(cmd) == 3 {
		reactToMessageId(cmd[1], cmd[2])
	} else if len(cmd) == 2 {
		reactToMessage(cmd[1])
	}

}

func reactToMessage(reaction string) {
	chat := k.NewChat(channel)
	chat.React(lastMessage.ID, reaction)

}
func reactToMessageId(messageId string, reaction string) {
	chat := k.NewChat(channel)
	ID, _ := strconv.Atoi(messageId)
	chat.React(ID, reaction)
}
