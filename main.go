package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"samhofi.us/x/keybase"
)

// Configurable section
var downloadPath = "/tmp/"
var outputFormat = "┌──[$USER@$DEVICE] [$ID] [$DATE - $TIME]\n└╼ $MSG"

// 02 = Day, Jan = Month, 06 = Year
var dateFormat = "02Jan06"

// 15 = hours, 04 = minutes, 05 = seconds
var timeFormat = "15:04"

// End configurable section

var k = keybase.NewKeybase()
var channel keybase.Channel
var channels []keybase.Channel
var stream = false
var lastMessage keybase.ChatAPI

func main() {
	if !k.LoggedIn {
		fmt.Println("You are not logged in.")
		return
	}

	kbtui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Printf("%+v", err)
	}
	defer kbtui.Close()
	kbtui.SetManagerFunc(layout)

	printToView(kbtui, "Chat", fmt.Sprintf("Welcome %s!", k.Username))
	go populateList(kbtui)
	go updateChatWindow(kbtui)
	if err := initKeybindings(kbtui); err != nil {
		log.Printf("%+v", err)
	}
	if err := kbtui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Printf("%+v", err)
	}
}
func populateChat(g *gocui.Gui) {
	lastMessage.ID = 0
	chat := k.NewChat(channel)
	maxX, _ := g.Size()
	api, err := chat.Read(maxX / 2)
	if err != nil {
		for _, testChan := range channels {
			if channel.Name == testChan.Name {
				channel = testChan
				channel.TopicName = "general"
			}
		}
		chat = k.NewChat(channel)
		_, err2 := chat.Read(2)
		if err2 != nil {
			printToView(g, "Feed", fmt.Sprintf("%+v", err))
			return
		} else {
			go populateChat(g)
			return
		}
	}
	var printMe []string
	var actuallyPrintMe string
	lastMessage.ID = api.Result.Messages[0].Msg.ID
	for _, message := range api.Result.Messages {
		if message.Msg.Content.Type == "text" {
			if lastMessage.ID < 1 {
				lastMessage.ID = message.Msg.ID
			}
			var apiCast keybase.ChatAPI
			apiCast.Msg = &message.Msg
			newMessage := formatOutput(apiCast)
			printMe = append(printMe, newMessage)
		}
	}
	for i := len(printMe) - 1; i >= 0; i-- {
		actuallyPrintMe += printMe[i]
		if i > 0 {
			actuallyPrintMe += "\n"
		}
	}
	printToView(g, "Chat", actuallyPrintMe)

}

func sendChat(message string, g *gocui.Gui) {
	chat := k.NewChat(channel)
	_, err := chat.Send(message)
	if err != nil {
		printToView(g, "Feed", fmt.Sprintf("There was an error %+v", err))
	}
}
func formatOutput(api keybase.ChatAPI) string {
	ret := ""
	if api.Msg.Content.Type == "text" {
		ret = outputFormat
		tm := time.Unix(int64(api.Msg.SentAt), 0)
		ret = strings.Replace(ret, "$MSG", api.Msg.Content.Text.Body, 1)
		ret = strings.Replace(ret, "$USER", api.Msg.Sender.Username, 1)
		ret = strings.Replace(ret, "$DEVICE", api.Msg.Sender.DeviceName, 1)
		ret = strings.Replace(ret, "$ID", fmt.Sprintf("%d", api.Msg.ID), 1)
		ret = strings.Replace(ret, "$DATE", fmt.Sprintf("%s", tm.Format(dateFormat)), 1)
		ret = strings.Replace(ret, "$TIME", fmt.Sprintf("%s", tm.Format(timeFormat)), 1)
	}
	return ret
}
func uploadFile(g *gocui.Gui, filePath string, fileName string) {
	chat := k.NewChat(channel)
	_, err := chat.Upload(fileName, filePath)
	if err != nil {
		printToView(g, "Feed", fmt.Sprintf("There was an error uploading %s to %s", filePath, channel.Name))
	} else {
		printToView(g, "Feed", fmt.Sprintf("Uploaded %s to %s", filePath, channel.Name))
	}
}
func downloadFile(g *gocui.Gui, messageID int, fileName string) {
	chat := k.NewChat(channel)
	_, err := chat.Download(messageID, fmt.Sprintf("%s/%s", downloadPath, fileName))
	if err != nil {
		printToView(g, "Feed", fmt.Sprintf("There was an error downloading %s from %s", fileName, channel.Name))
	} else {
		printToView(g, "Feed", fmt.Sprintf("Downloaded %s from %s", fileName, channel.Name))
	}
}

func populateList(g *gocui.Gui) {
	_, maxY := g.Size()
	if testVar, err := k.ChatList(); err != nil {
		log.Printf("%+v", err)
	} else {

		clearView(g, "List")
		var recentPMs = "---[PMs]---\n"
		var recentPMsCount = 0
		var recentChannels = "---[Teams]---\n"
		var recentChannelsCount = 0
		for _, s := range testVar.Result.Conversations {
			channels = append(channels, s.Channel)
			if s.Channel.MembersType == keybase.TEAM {
				recentChannelsCount++
				if recentChannelsCount <= ((maxY - 2) / 3) {
					if s.Unread {
						recentChannels += "*"
					}
					recentChannels += fmt.Sprintf("%s\n\t#%s\n", s.Channel.Name, s.Channel.TopicName)
				}
			} else {
				recentPMsCount++
				if recentPMsCount <= ((maxY - 2) / 3) {
					if s.Unread {
						recentPMs += "*"
					}
					recentPMs += fmt.Sprintf("%s\n", cleanChannelName(s.Channel.Name))
				}
			}
		}
		time.Sleep(1 * time.Millisecond)
		printToView(g, "List", fmt.Sprintf("%s%s", recentPMs, recentChannels))
	}
}

func clearView(kbtui *gocui.Gui, viewName string) {
	kbtui.Update(func(g *gocui.Gui) error {
		inputView, err := kbtui.View(viewName)
		if err != nil {
			return err
		} else {
			inputView.Clear()
			inputView.SetCursor(0, 0)
			inputView.SetOrigin(0, 0)
		}
		return nil
	})

}

func printToView(kbtui *gocui.Gui, viewName string, message string) {
	kbtui.Update(func(g *gocui.Gui) error {
		updatingView, err := kbtui.View(viewName)
		if err != nil {
			return err
		} else {
			fmt.Fprintf(updatingView, "%s\n", message)
		}
		return nil
	})
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if feedView, err := g.SetView("Feed", maxX/2-maxX/3, 0, maxX-1, maxY/5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		feedView.Autoscroll = true
		feedView.Wrap = true
		fmt.Fprintln(feedView, "Feed Window - If you are mentioned or receive a PM it will show here")
	}
	if chatView, err2 := g.SetView("Chat", maxX/2-maxX/3, maxY/5+1, maxX-1, maxY-5); err2 != nil {
		if err2 != gocui.ErrUnknownView {
			return err2
		}
		chatView.Autoscroll = true
		chatView.Wrap = true
		fmt.Fprintf(chatView, "Your chats will appear here.\nSupported commands are as follows:\n")
		fmt.Fprintln(chatView, "/j $username - Open your chat with $username")
		fmt.Fprintln(chatView, "/j $team $channel - Open $channel from $team")
		fmt.Fprintln(chatView, "/u $path $title - Uploads file $path with title $title")
		fmt.Fprintln(chatView, "/d $msgId $downloadName - Downloads file from $msgId to $DownloadPath/$downloadName")
		fmt.Fprintln(chatView, "/r $msgId $reaction - Reacts to $msgId with $reaction reaction can be emoji :+1:")
		fmt.Fprintln(chatView, "      Can also be used for STRING reactions")
		fmt.Fprintln(chatView, "/s  - Experimental: View all incoming messages from everywhere.")
		fmt.Fprintln(chatView, "/q - Exit")
	}
	if inputView, err3 := g.SetView("Input", maxX/2-maxX/3, maxY-4, maxX-1, maxY-1); err3 != nil {
		if err3 != gocui.ErrUnknownView {
			return err3
		}
		if _, err := g.SetCurrentView("Input"); err != nil {
			return err
		}
		inputView.Editable = true
		inputView.Wrap = true
		g.Cursor = true
	}
	if listView, err4 := g.SetView("List", 0, 0, maxX/2-maxX/3-1, maxY-1); err4 != nil {
		if err4 != gocui.ErrUnknownView {
			return err4
		}
		fmt.Fprintf(listView, "Lists\nWindow\nTo view\n activity")
	}
	return nil
}

func getInputString(g *gocui.Gui) (string, error) {
	inputView, _ := g.View("Input")
	return inputView.Line(0)
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			input, err := getInputString(g)
			if err != nil {
				return err
			}
			if input != "" {
				clearView(g, "Input")
				return nil
			} else {
				return gocui.ErrQuit
			}
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Input", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return handleInput(g)
		}); err != nil {
		return err
	}
	return nil
}

func updateChatWindow(g *gocui.Gui) {
	k.Run(func(api keybase.ChatAPI) {
		handleMessage(api, g)
	})

}

func cleanChannelName(c string) string {
	newChannelName := strings.Replace(c, fmt.Sprintf("%s,", k.Username), "", 1)
	return strings.Replace(newChannelName, fmt.Sprintf(",%s", k.Username), "", 1)
}

func handleMessage(api keybase.ChatAPI, g *gocui.Gui) {
	if api.Msg.Content.Type == "text" {
		go populateList(g)
		msgBody := api.Msg.Content.Text.Body
		msgSender := api.Msg.Sender.Username
		channelName := api.Msg.Channel.Name
		if !stream {
			if msgSender != k.Username {
				if api.Msg.Channel.MembersType == keybase.TEAM {
					topicName := api.Msg.Channel.TopicName
					for _, m := range api.Msg.Content.Text.UserMentions {
						if m.Text == k.Username {
							// We are in a team
							if topicName != channel.TopicName {
								printToView(g, "Feed", fmt.Sprintf("[ %s#%s ] %s: %s", channelName, topicName, msgSender, msgBody))
							}

							break
						}
					}
				} else {
					if msgSender != channel.Name {
						printToView(g, "Feed", fmt.Sprintf("PM from @%s: %s", cleanChannelName(channelName), msgBody))
					}

				}
			}
			if api.Msg.Channel.MembersType == channel.MembersType && cleanChannelName(api.Msg.Channel.Name) == channel.Name {
				if channel.MembersType == keybase.TEAM && channel.TopicName != api.Msg.Channel.TopicName {
					// Do nothing, wrong channel
				} else {

					printToView(g, "Chat", formatOutput(api))
					chat := k.NewChat(channel)
					lastMessage.ID = api.Msg.ID
					chat.Read(api.Msg.ID)
				}

			}
		} else {
			if api.Msg.Channel.MembersType == keybase.TEAM {
				topicName := api.Msg.Channel.TopicName
				printToView(g, "Chat", fmt.Sprintf("@%s#%s [%s]: %s", channelName, topicName, msgSender, msgBody))
			} else {
				printToView(g, "Chat", fmt.Sprintf("PM @%s [%s]: %s", cleanChannelName(channelName), msgSender, msgBody))
			}
		}
	} else {
		//TODO: For edit/delete run this
		if api.Msg.Channel.MembersType == channel.MembersType && cleanChannelName(api.Msg.Channel.Name) == channel.Name {
			go populateChat(g)
		}
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
func handleInput(g *gocui.Gui) error {
	inputString, _ := getInputString(g)
	if inputString == "" {
		return nil
	}
	command := strings.Split(inputString, " ")

	switch strings.ToLower(command[0]) {
	case "/q":
		return gocui.ErrQuit
	case "/j":
		stream = false
		if len(command) == 3 {
			channel.MembersType = keybase.TEAM
			channel.Name = command[1]
			channel.TopicName = command[2]
			printToView(g, "Feed", fmt.Sprintf("You are joining: @%s#%s", channel.Name, channel.TopicName))
			clearView(g, "Chat")
			go populateChat(g)
		} else if len(command) == 2 {
			channel.MembersType = keybase.USER
			channel.Name = command[1]
			channel.TopicName = ""
			printToView(g, "Feed", fmt.Sprintf("You are joining: @%s", channel.Name))
			clearView(g, "Chat")
			go populateChat(g)
		} else {
			printToView(g, "Feed", "To join a team use /j <team> <channel>")
			printToView(g, "Feed", "To join a PM use /j <user>")
		}
	case "/u":
		if len(command) == 3 {
			filePath := command[1]
			fileName := command[2]
			uploadFile(g, filePath, fileName)
		} else if len(command) == 2 {
			filePath := command[1]
			fileName := ""
			uploadFile(g, filePath, fileName)
		} else {
			printToView(g, "Feed", "To upload a file, supply full path and optional title (no spaces)")
		}
	case "/d":
		if len(command) == 3 {
			messageId, err := strconv.Atoi(command[1])
			if err != nil {
				printToView(g, "Feed", "Invalid message ID")
			} else {
				fileName := command[2]
				downloadFile(g, messageId, fileName)
			}
		} else if len(command) == 2 {
			messageId, err := strconv.Atoi(command[1])
			if err != nil {
				printToView(g, "Feed", "Invalid message ID")
			} else {
				downloadFile(g, messageId, command[1])
			}
		}
	case "/s":
		clearView(g, "Chat")
		stream = true
		printToView(g, "Feed", "You have begun viewing the formatted stream.")
	case "/r":
		if len(command) == 3 {
			reactToMessageId(command[1], command[2])
		} else {
			printToView(g, "Feed", "/r $messageId $desiredReaction")
		}
	default:
		if inputString[:1] == "+" {
			reactToMessage(strings.Replace(inputString, "+", "", 1))
		} else {
			go sendChat(inputString, g)
		}
		go populateList(g)
	}
	clearView(g, "Input")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
