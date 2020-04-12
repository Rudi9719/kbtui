// +build !rm_basic_commands allcommands devcmd

package main

import (
	"fmt"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"dev"},
		Description: "- Switch to dev channels",
		Help:        "",
		Exec:        cmdDev,
	}

	RegisterCommand(command)
}

func cmdDev(cmd []string) {
	dev = !dev

	printInfo(fmt.Sprintf("You have toggled the dev flag to %+v", dev))

  if (lastChat != "") {
    // Switching from regular to dev mode? Dev chats don't use channels. Let's strip the channel name.
    n := ""
    if (dev) { n = strings.Split(lastChat, "#")[0] } else { n = lastChat }
    cmdJoin([]string{"/join", n})
  }
  //go updateChatWindow() // Otherwise you won't be able to process incoming messages.
}
