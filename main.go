package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"samhofi.us/x/keybase"
	"strings"
	"time"
)

var k = keybase.NewKeybase()
var teamOrUser = "home"
var channel = ""
var myUsername = ""

func main() {
	if k.LoggedIn == true {
		kbtui, err := gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			log.Panicln(err)
		}
		defer kbtui.Close()

		kbtui.SetManagerFunc(layout)
		go loginGreeter(kbtui)
		go populateList(kbtui)
		go updateChatWindow(kbtui)
		if err := initKeybindings(kbtui); err != nil {
			log.Fatalln(err)
		}
		if err := kbtui.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	} else {
		fmt.Println("You are not logged in.")
		return
	}

}
func loginGreeter(g *gocui.Gui) {
	myUsername = k.Username
	messg := "Welcome " + myUsername+ "!"
	printToView(g, "Chat", messg)
}
func sendToUser(msg string) {
	chat := k.NewChat(keybase.Channel{Name:teamOrUser})
	chat.Send(msg)
}

func sendToTeam(msg string) {
	chat := k.NewChat(keybase.Channel{
		Name:        teamOrUser,
		MembersType: "team",
		TopicName:   channel,
	})
	chat.Send(msg)
}
func populateList(g *gocui.Gui) {
	for {
		if testVar, err := k.ChatList(); err != nil {
			log.Fatalln(err)
		} else {
			clearView(g, "List")
			for _, s := range testVar.Result.Conversations {
				printToView(g, "List", s.Channel.Name)
			}
		}
		time.Sleep(5 * time.Second)
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
			fmt.Fprintf(updatingView, message + "\n")
		}
		return nil
	})
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if feedView, err := g.SetView("Feed", 12, 0, maxX-1, maxY/5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		feedView.Autoscroll = true
		feedView.Wrap = true
		fmt.Fprintln(feedView, "Feed Window - If you are mentioned or receive a PM it will show here")
	}
	if chatView, err2 := g.SetView("Chat", 12, maxY/5+1, maxX-1, maxY-5); err2 != nil {
		if err2 != gocui.ErrUnknownView {
			return err2
		}
		chatView.Autoscroll = true
		chatView.Wrap = true
		fmt.Fprintf(chatView, "Your chats will appear here.\nSupported commands are as follows:\n")
		fmt.Fprintln(chatView, "/j $username - Open your chat with $username")
		fmt.Fprintln(chatView, "/j $team $channel - Open $channel from $team")
		fmt.Fprintln(chatView, "          Please note: small teams only have #general")
		fmt.Fprintln(chatView, "/q - Exit")
	}
	if inputView, err3 := g.SetView("Input", 12, maxY-4, maxX-1, maxY-1); err3 != nil {
		if err3 != gocui.ErrUnknownView {
			return err3
		}
		 if _, err := g.SetCurrentView("Input"); err != nil {
		 	return err
		 }
		inputView.Editable = true
		inputView.Wrap = true
	}
	if listView, err4 := g.SetView("List", 0, 0, 10, maxY-1); err4 != nil {
		if err4 != gocui.ErrUnknownView {
			return err4
		}
		listView.Wrap = true
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

func handleMessage(api keybase.ChatAPI, g *gocui.Gui) {
	if api.Msg.Content.Type == "text" {
		if strings.Contains(api.Msg.Content.Text.Body, myUsername) || strings.Contains(api.Msg.Channel.Name, myUsername) {
			if api.Msg.Sender.Username != myUsername {
				message := api.Msg.Content.Text.Body
				sender := api.Msg.Sender.Username
				team := api.Msg.Channel.Name
				channel := api.Msg.Channel.TopicName
				printMe := "@" + team + "#" + channel + " " + sender + ": " + message
				printToView(g, "Feed", printMe)
			}
		}
		if strings.Contains(api.Msg.Channel.Name, teamOrUser) {
			if channel != "" && api.Msg.Channel.TopicName == channel {

			} else {
				message := api.Msg.Content.Text.Body
				sender := api.Msg.Sender.Username
				printToView(g, "Chat", sender+": "+message)
			}
		}
	}
}
func handleInput(g *gocui.Gui) error {
	inputString, _ := getInputString(g)
	command := ""
	if len(inputString) > 2 {
		command = inputString[:2]
	}
	if "/q" == command {
		return gocui.ErrQuit
	} else if "/j" == command {
		clearView(g, "Chat")
		arr := strings.Split(inputString, " ")
		if len(arr) > 2 {
			teamOrUser = arr[1]
			channel = arr[2]
			printToView(g, "Feed", "You have joined: @" + teamOrUser + "#" + channel)
		} else {
			teamOrUser = arr[1]
			channel = ""
			printToView(g, "Feed", "You have joined: @" + teamOrUser)
		}
	} else {
		if channel == "" {
			sendToUser(inputString)
		} else {
			sendToTeam(inputString)
		}
	}
	clearView(g, "Input")
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
