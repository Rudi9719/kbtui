// +ignore
// +build allcommands postcmd

package main

import (
	"fmt"
	"strings"

	"samhofi.us/x/keybase"
)

func init() {
	command := Command{
		Cmd:         []string{"post"},
		Description: "- Post public messages on your wall",
		Help:        "",
		Exec:        cmdPost,
	}

	RegisterCommand(command)
}
func cmdPost(cmd []string) {
	var pubChan keybase.Channel
	pubChan.Public = true
	pubChan.MembersType = keybase.USER
	pubChan.Name = k.Username
	post := strings.Join(cmd[1:], " ")
	chat := k.NewChat(pubChan)
	_, err := chat.Send(post)
	if err != nil {
		printError(fmt.Sprintf("There was an error with your post: %+v", err))
	} else {
		printInfo("You have publically posted to your wall, signed by your current device.")
	}
}
