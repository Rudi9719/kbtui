// +ignore
// +build allcommands shrugcmd

package main

import "strings"

func init() {
	command := Command{
		Cmd:         []string{"shrug", "shrg"},
		Description: "$message - append a shrug ( ¯\\_(ツ)_/¯ )to your message",
		Help:        "",
		Exec:        cmdShrug,
	}

	RegisterCommand(command)
}

func cmdShrug(cmd []string) {
	cmd = append(cmd, " ¯\\_(ツ)_/¯")

	sendChat(strings.Join(cmd[1:], " "))
}
