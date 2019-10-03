package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"samhofi.us/x/keybase"
)

var commands = make(map[string]Command)
var baseCommands = make([]string, 0)

var k = keybase.NewKeybase()
var channel keybase.Channel
var channels []keybase.Channel
var stream = false
var lastMessage keybase.ChatAPI
var g *gocui.Gui

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
	g = kbtui
	go populateList()
	go updateChatWindow()
	if err := initKeybindings(); err != nil {
		log.Printf("%+v", err)
	}
	if err := kbtui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Printf("%+v", err)
	}
}
func populateChat() {
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
			printToView("Feed", fmt.Sprintf("%+v", err))
			return
		} else {
			go populateChat()
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
	printToView("Chat", actuallyPrintMe)

}

func sendChat(message string) {
	chat := k.NewChat(channel)
	_, err := chat.Send(message)
	if err != nil {
		printToView("Feed", fmt.Sprintf("There was an error %+v", err))
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

func populateList() {
	_, maxY := g.Size()
	if testVar, err := k.ChatList(); err != nil {
		log.Printf("%+v", err)
	} else {

		clearView("List")
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
		printToView("List", fmt.Sprintf("%s%s", recentPMs, recentChannels))
	}
}

func clearView(viewName string) {
	g.Update(func(g *gocui.Gui) error {
		inputView, err := g.View(viewName)
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

func printToView(viewName string, message string) {
	g.Update(func(g *gocui.Gui) error {
		updatingView, err := g.View(viewName)
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
		fmt.Fprintf(chatView, "Welcome %s!\n\nYour chats will appear here.\nSupported commands are as follows:\n\n", k.Username)
		RunCommand("help")
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

func getInputString() (string, error) {
	inputView, _ := g.View("Input")
	return inputView.Line(0)
}

func initKeybindings() error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			input, err := getInputString()
			if err != nil {
				return err
			}
			if input != "" {
				clearView("Input")
				return nil
			} else {
				return gocui.ErrQuit
			}
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Input", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return handleInput()
		}); err != nil {
		return err
	}
	return nil
}

func updateChatWindow() {
	k.Run(func(api keybase.ChatAPI) {
		handleMessage(api)
	})

}

func cleanChannelName(c string) string {
	newChannelName := strings.Replace(c, fmt.Sprintf("%s,", k.Username), "", 1)
	return strings.Replace(newChannelName, fmt.Sprintf(",%s", k.Username), "", 1)
}

func handleMessage(api keybase.ChatAPI) {
	if api.Msg.Content.Type == "text" {
		go populateList()
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
								printToView("Feed", fmt.Sprintf("[ %s#%s ] %s: %s", channelName, topicName, msgSender, msgBody))
							}

							break
						}
					}
				} else {
					if msgSender != channel.Name {
						printToView("Feed", fmt.Sprintf("PM from @%s: %s", cleanChannelName(channelName), msgBody))
					}

				}
			}
			if api.Msg.Channel.MembersType == channel.MembersType && cleanChannelName(api.Msg.Channel.Name) == channel.Name {
				if channel.MembersType == keybase.TEAM && channel.TopicName != api.Msg.Channel.TopicName {
					// Do nothing, wrong channel
				} else {

					printToView("Chat", formatOutput(api))
					chat := k.NewChat(channel)
					lastMessage.ID = api.Msg.ID
					chat.Read(api.Msg.ID)
				}

			}
		} else {
			if api.Msg.Channel.MembersType == keybase.TEAM {
				topicName := api.Msg.Channel.TopicName
				printToView("Chat", fmt.Sprintf("@%s#%s [%s]: %s", channelName, topicName, msgSender, msgBody))
			} else {
				printToView("Chat", fmt.Sprintf("PM @%s [%s]: %s", cleanChannelName(channelName), msgSender, msgBody))
			}
		}
	} else {
		//TODO: For edit/delete run this
		if api.Msg.Channel.MembersType == channel.MembersType && cleanChannelName(api.Msg.Channel.Name) == channel.Name {
			go populateChat()
		}
	}
}

func handleInput() error {
	clearView("Input")
	inputString, _ := getInputString()
	if inputString == "" {
		return nil
	}
	if strings.HasPrefix(inputString, cmdPrefix) {
		cmd := strings.Split(inputString[len(cmdPrefix):], " ")
		if c, ok := commands[cmd[0]]; ok {
			c.Exec(cmd)
			return nil
		} else if cmd[0] == "q" || cmd[0] == "quit" {
			return gocui.ErrQuit
		} else {
			printToView("Feed", fmt.Sprintf("Command '%s' not recognized", cmd[0]))
			return nil
		}
	}
	if inputString[:1] == "+" {
		cmd := strings.Split(inputString, " ")
		RunCommand(cmd...)
	} else {
		go sendChat(inputString)
	}

	go populateList()
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// RegisterCommand registers a command to be used within the client
func RegisterCommand(c Command) error {
	var notAdded string
	for i, cmd := range c.Cmd {
		if _, ok := commands[cmd]; !ok {
			if i == 0 {
				baseCommands = append(baseCommands, cmd)
			}
			commands[cmd] = c
			continue
		}
		notAdded = fmt.Sprintf("%s, %s", notAdded, cmd)
	}
	if notAdded != "" {
		return fmt.Errorf("The following aliases were not added because they already exist: %s", notAdded)
	}
	return nil
}

// RunCommand calls a command as if it was run by the user
func RunCommand(c ...string) {
	commands[c[0]].Exec(c)
}
