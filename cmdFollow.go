// +build !rm_basic_commands allcommands followcmd

package main

import (
	"fmt"
)

func init() {
	command := Command{
		Cmd:         []string{"follow"},
		Description: "$username - Follows the given user",
		Help:        "",
		Exec:        cmdFollow,
	}
	RegisterCommand(command)
}

func cmdFollow(cmd []string) {
	if len(cmd) == 2 {
		go follow(cmd[1])
	} else {
		printFollowHelp()
	}
}
func follow(username string) {
	k.Exec("follow", username, "-y")
	printInfoF("Now follows $TEXT", config.Colors.Message.LinkKeybase.stylize(username))
	followedInSteps[username] = 1
}

func printFollowHelp() {
	printInfo(fmt.Sprintf("To follow a user use %sfollow <username>", config.Basics.CmdPrefix))
}
