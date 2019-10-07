// +build !rm_basic_commands allcommands wallcmd

package main

import (
	"fmt"
	"strings"
	"time"

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
	go cmdPopulateWall(cmd)
}

func cmdPopulateWall(cmd []string) {
	var users []keybase.Channel
	var requestedUsers string
	var printMe []string
	var actuallyPrintMe string
	start := time.Now()
	if len(cmd) > 1 {
		if cmd[1] == "!all" {
			go cmdAllWall()
			return
		}
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
		printToView("Feed", fmt.Sprintf("Error, can't run wall in teams"))
		go cmdAllWall()
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
			if len(users) < 6 {
				printToView("Feed", fmt.Sprintf("There was an error for user %s: %+v", cleanChannelName(chann.Name), err))
			}
		} else {
			for _, message := range api.Result.Messages {
				if message.Msg.Content.Type == "text" {
					var apiCast keybase.ChatAPI
					apiCast.Msg = &message.Msg
					newMessage := formatOutput(apiCast)
					printMe = append(printMe, newMessage)
				}
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
	time.Sleep(1 * time.Millisecond)
	printToView("Chat", fmt.Sprintf("Your wall query took %s", time.Since(start)))
}
func cmdAllWall() {
	bytes, _ := k.Exec("list-following")
	bigString := string(bytes)
	following := strings.Split(bigString, "\n")
	go cmdPopulateWall(following)
}
