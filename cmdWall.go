// +ignore
// +build allcommands wallcmd

package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"samhofi.us/x/keybase"
)

func init() {
	command := Command{
		Cmd:         []string{"wall", "w"},
		Description: "$user / !all - Show public messages for a user or all users you follow",
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
	result := make(map[int]keybase.ChatAPI)
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
		requestedUsers += k.Username
		var newChan keybase.Channel
		newChan.MembersType = keybase.USER
		newChan.Name = k.Username
		newChan.TopicName = ""
		newChan.Public = true
		users = append(users, newChan)
	}
	if len(users) < 1 {
		return
	}

	printInfoF("Displaying public messages for user $TEXT", config.Colors.Message.LinkKeybase.stylize(requestedUsers))
	for _, chann := range users {
		chat := k.NewChat(chann)
		api, err := chat.Read()
		if err != nil {
			if len(users) < 6 {
				printError(fmt.Sprintf("There was an error for user %s: %+v", cleanChannelName(chann.Name), err))
				return
			}
		} else {
			for i, message := range api.Result.Messages {
				if message.Msg.Content.Type == "text" {
					var apiCast keybase.ChatAPI
					apiCast.Msg = &api.Result.Messages[i].Msg
					result[apiCast.Msg.SentAt] = apiCast
					newMessage := formatOutput(apiCast)
					printMe = append(printMe, newMessage)

				}
			}
		}

	}

	keys := make([]int, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	time.Sleep(1 * time.Millisecond)
	for _, k := range keys {
		actuallyPrintMe += formatOutput(result[k]) + "\n"
	}
	printToView("Chat", fmt.Sprintf("\n<Wall>\n\n%s\nYour wall query took %s\n</Wall>\n", actuallyPrintMe, time.Since(start)))
}
func cmdAllWall() {
	bytes, _ := k.Exec("list-following")
	bigString := string(bytes)
	following := strings.Split(bigString, "\n")
	go cmdPopulateWall(following)
}
