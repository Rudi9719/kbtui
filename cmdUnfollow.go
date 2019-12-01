// +build !rm_basic_commands allcommands followcmd

package main

import (
	"fmt"
)

func init() {
	command := Command{
		Cmd:         []string{"unfollow"},
		Description: "$username - Unfollows the given user",
		Help:        "",
		Exec:        cmdUnfollow,
	}
	RegisterCommand(command)
}

func cmdUnfollow(cmd []string) {
	if len(cmd) == 2 {
		go unfollow(cmd[1])
	} else {
		printUnfollowHelp()
	}
}
func unfollow(username string) {
	k.Exec("unfollow", username)
	printInfoF("Now unfollows $TEXT", config.Colors.Message.LinkKeybase.stylize(username))
}

func printUnfollowHelp() {
	printInfo(fmt.Sprintf("To unfollow a user use %sunfollow <username>", config.Basics.CmdPrefix))
}
