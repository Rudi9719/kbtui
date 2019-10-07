// +build !rm_basic_commands allcommands postcmd

package main

import (

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
	chat.Send(post)
}
