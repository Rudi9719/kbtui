package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"samhofi.us/x/keybase"
)

var k = keybase.NewKeybase()
var channel keybase.Channel
var channels [] keybase.Channel
var stream bool = false
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
	chat := k.NewChat(channel)
	maxX, _ := g.Size()
	if api, err := chat.Read(maxX/2); err != nil {
		log.Printf("%+v", err)
	} else {
		var printMe []string
		var actuallyPrintMe string
		firstmsg := true
		for _, message := range api.Result.Messages {
			if message.Msg.Content.Type == "text" {
				if (firstmsg) {
					firstmsg = false
					lastMessage.ID = message.Msg.ID
				}
				msgSender := message.Msg.Sender.Username
				msgBody := message.Msg.Content.Text.Body
				newMessage := fmt.Sprintf("[%s]: %s", msgSender, msgBody)
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
}

func sendChat(message string) {
	chat := k.NewChat(channel)
	chat.Send(message)
}

func populateList(g *gocui.Gui) {
	_, maxY := g.Size()
	if testVar, err := k.ChatList(); err != nil {
		log.Printf("%+v",err)
	} else {
		clearView(g, "List")
		var recentPMs = "---[PMs]---\n"
		var recentPMsCount = 0
		var recentChannels = "---[Teams]---\n"
		var recentChannelsCount = 0
		for _, s := range testVar.Result.Conversations {
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
				if recentPMsCount <= ((maxY- 2) / 3) {
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
		fmt.Fprintln(chatView, "/s  - Experimental: View all incoming messages from everywhere.")
		fmt.Fprintln(chatView, "          Please note: small teams only have #general")
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
			return gocui.ErrQuit
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
							printToView(g, "Feed", fmt.Sprintf("[ %s#%s ] %s: %s", channelName, topicName, msgSender, msgBody))
							break
						}
					}
				} else {
					printToView(g, "Feed", fmt.Sprintf("PM from @%s: %s", cleanChannelName(channelName), msgBody))
				}
			}
			if api.Msg.Channel.MembersType == channel.MembersType && cleanChannelName(api.Msg.Channel.Name) == channel.Name {
				printToView(g, "Chat", fmt.Sprintf("[%s]: %s", msgSender, msgBody))
				chat := k.NewChat(channel)
				chat.Read(api.Msg.ID)
			}
		} else {
			if api.Msg.Channel.MembersType == keybase.TEAM {
				topicName := api.Msg.Channel.TopicName
				printToView(g, "Chat", fmt.Sprintf("@%s#%s [%s]: %s", channelName, topicName, msgSender, msgBody))
			} else {
				printToView(g, "Chat", fmt.Sprintf("PM @%s [%s]: %s", cleanChannelName(channelName), msgSender, msgBody))
			}
		}
		lastMessage = api
	}
}
func reactToMessage(reaction string) {
	chat := k.NewChat(channel)
	chat.React(lastMessage.ID, reaction)
}
func handleInput(g *gocui.Gui) error {
	inputString, _ := getInputString(g)
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
			printToView(g, "Feed", fmt.Sprintf("You have joined: @%s#%s", channel.Name, channel.TopicName))
			clearView(g, "Chat")
			go populateChat(g)
		} else if len(command) == 2 {
			channel.MembersType = keybase.USER
			channel.Name = command[1]
			channel.TopicName = ""
			printToView(g, "Feed", fmt.Sprintf("You have joined: @%s", channel.Name))
			clearView(g, "Chat")
			go populateChat(g)
		} else {
			printToView(g, "Feed", "To join a team use /j <team> <channel>")
			printToView(g, "Feed", "To join a PM use /j <user>")
		}
	case "/s":
		clearView(g, "Chat")
		stream = true
		printToView(g, "Feed", "You have begun viewing the formatted stream.")
	default:
		if inputString[:1] == "+" {
			reactToMessage(strings.Replace(inputString, "+", "", 1))
		} else {
			go sendChat(inputString)
		}
	}
	clearView(g, "Input")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
