// +build !rm_basic_commands allcommands wallcmd

package main

import (
	"fmt"

	"samhofi.us/x/keybase"
)

func init() {
	command := Command{
		Cmd:         []string{"wall", "w"},
		Description: "- Show public messages for a user",
		Help:        "",
		Exec:        cmdWall,
	}

	RegisterCommand(command)
}

func cmdWall(cmd []string) {
	var users []keybase.Channel
	var requestedUsers string
	var printMe []string
	var actuallyPrintMe string
	if len(cmd) > 1 {
		for _, username := range cmd[1:] {
			requestedUsers += fmt.Sprintf("%s ", username)
			var newChan keybase.Channel
			newChan.MembersType = keybase.USER
			newChan.Name = username
			newChan.TopicName = ""
			newChan.Public = true
			users = append(users, newChan)
		}
	} else if channel.MembersType == keybase.USER {
		users = append(users, channel)
		users[0].Public = true
		requestedUsers += cleanChannelName(channel.Name)

	} else {
		printToView("Feed", fmt.Sprintf("%+v", "\nError", channel.MembersType))
		return
	}
	if len(users) < 1 {
		return
	}
	printToView("Feed", fmt.Sprintf("Displaying public messages for user %s", requestedUsers))
	for _, chann := range users {
		chat := k.NewChat(chann)
		api, err := chat.Read()
		if err != nil {
			printToView("Feed", fmt.Sprintf("There was an error for user %s: %+v", cleanChannelName(chann.Name), err))
			return
		}
		for _, message := range api.Result.Messages {
			if message.Msg.Content.Type == "text" {
				var apiCast keybase.ChatAPI
				apiCast.Msg = &message.Msg
				newMessage := formatOutput(apiCast)
				printMe = append(printMe, newMessage)
			}
		}

	}
	for i := len(printMe) - 1; i >= 0; i-- {
		actuallyPrintMe += printMe[i]
		if i > 0 {
			actuallyPrintMe += "\n"
		}
	}
	printToView("Chat", fmt.Sprintf("\nWall:\n%s\n", actuallyPrintMe))
}
