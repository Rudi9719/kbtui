package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
	"samhofi.us/x/keybase"
)

var (
	typeCommands = make(map[string]TypeCommand)
	commands     = make(map[string]Command)
	baseCommands = make([]string, 0)

	dev         = false
	k           = keybase.NewKeybase()
	channel     keybase.Channel
	channels    []keybase.Channel
	stream      = false
	lastMessage keybase.ChatAPI
	lastChat    = ""
	g           *gocui.Gui
)

var config *Config

func main() {
	if !k.LoggedIn {
		fmt.Println("You are not logged in.")
		return
	}
	var err error
	g, err = gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)
	RunCommand("config", "load")
	go populateList()
	go updateChatWindow()
	if len(os.Args) > 1 {
		os.Args[0] = "join"
		RunCommand(os.Args...)

	}
	fmt.Println("initKeybindings")
	if err := initKeybindings(); err != nil {
		fmt.Printf("%+v", err)
	}
	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		fmt.Printf("%+v", err)
	}
	go generateChannelTabCompletionSlice()
}

// Gocui basic setup
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if editView, err := g.SetView("Edit", maxX/2-maxX/3+1, maxY/2, maxX-2, maxY/2+10, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		editView.Editable = true
		editView.Wrap = true
		fmt.Fprintln(editView, "Edit window. Should disappear")
	}
	if feedView, err := g.SetView("Feed", maxX/2-maxX/3, 0, maxX-1, maxY/5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		feedView.Autoscroll = true
		feedView.Wrap = true
		feedView.Title = "Feed Window"
		printInfo("Feed Window - If you are mentioned or receive a PM it will show here")
	}
	if chatView, err2 := g.SetView("Chat", maxX/2-maxX/3, maxY/5+1, maxX-1, maxY-5, 0); err2 != nil {
		if !gocui.IsUnknownView(err2) {
			return err2
		}
		chatView.Autoscroll = true
		chatView.Wrap = true
		welcomeText := basicStyle.stylize("Welcome $USER!\n\nYour chats will appear here.\nSupported commands are as follows:\n")
		welcomeText = welcomeText.replace("$USER", config.Colors.Message.Mention.stylize(k.Username))
		fmt.Fprintln(chatView, welcomeText.string())
		RunCommand("help")
	}
	if inputView, err3 := g.SetView("Input", maxX/2-maxX/3, maxY-4, maxX-1, maxY-1, 0); err3 != nil {
		if !gocui.IsUnknownView(err3) {
			return err3
		}
		if _, err := g.SetCurrentView("Input"); err != nil {
			return err
		}
		inputView.Editable = true
		inputView.Wrap = true
		inputView.Title = fmt.Sprintf(" Not in a chat - write `%sj` to join", config.Basics.CmdPrefix)
		g.Cursor = true
	}
	if listView, err4 := g.SetView("List", 0, 0, maxX/2-maxX/3-1, maxY-1, 0); err4 != nil {
		if !gocui.IsUnknownView(err4) {
			return err4
		}
		listView.Title = "Channels"
		fmt.Fprintf(listView, "Lists\nWindow\nTo view\n activity")
	}
	return nil
}
func scrollUP(v *gocui.View) error {
	scrollView(v, -1)
	return nil
}
func scrollDOWN(v *gocui.View) error {
	scrollView(v, 1)
	return nil
}
func scrollView(v *gocui.View, delta int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+delta); err != nil {
			return err
		}
	}
	return nil
}
func autoScrollView(vn string) error {
	v, err := g.View(vn)
	if err != nil {
		return err
	} else if v != nil {
		v.Autoscroll = true
	}
	return nil
}
func initKeybindings() error {
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cv, _ := g.View("Chat")
			err := scrollUP(cv)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cv, _ := g.View("Chat")
			err := scrollDOWN(cv)
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			input, err := getInputString("Input")
			if err != nil {
				return err
			}
			if input != "" {
				clearView("Input")
				return nil
			}
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlZ, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			cmdJoin([]string{"/join", lastChat})
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Edit", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			popupView("Chat")
			popupView("Input")
			clearView("Edit")
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Input", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return handleInput("Input")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Input", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return handleTab("Input")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Edit", gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			popupView("Chat")
			popupView("Input")
			return handleInput("Edit")

		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Input", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			RunCommand("edit")
			return nil
		}); err != nil {
		return err
	}
	return nil
}

// End gocui basic setup

// Gocui helper funcs
func setViewTitle(viewName string, title string) {
	g.Update(func(g *gocui.Gui) error {
		updatingView, err := g.View(viewName)
		if err != nil {
			return err
		}
		updatingView.Title = title
		return nil
	})
}
func getViewTitle(viewName string) string {
	view, err := g.View(viewName)
	if err != nil {
		// in case there is active tab completion, filter that to just the view title and not the completion options.
		printError(fmt.Sprintf("Error getting view title: %s", err))
		return ""
	}
	return strings.Split(view.Title, "||")[0]
}
func popupView(viewName string) {
	_, err := g.SetCurrentView(viewName)
	if err != nil {
		printError(fmt.Sprintf("%+v", err))
	}
	_, err = g.SetViewOnTop(viewName)
	if err != nil {
		printError(fmt.Sprintf("%+v", err))
	}
	g.Update(func(g *gocui.Gui) error {
		updatingView, err := g.View(viewName)
		if err != nil {
			return err
		}
		updatingView.MoveCursor(0, 0, true)

		return nil

	})
}
func moveCursorToEnd(viewName string) {
	g.Update(func(g *gocui.Gui) error {
		inputView, err := g.View(viewName)
		if err != nil {
			return err
		}
		inputString, _ := getInputString(viewName)
		stringLen := len(inputString)
		maxX, _ := inputView.Size()
		x := stringLen % maxX
		y := stringLen / maxX
		inputView.SetCursor(0, 0)
		inputView.SetOrigin(0, 0)
		inputView.MoveCursor(x, y, true)
		return nil

	})
}
func clearView(viewName string) {
	g.Update(func(g *gocui.Gui) error {
		inputView, err := g.View(viewName)
		if err != nil {
			return err
		}
		inputView.Clear()
		inputView.SetCursor(0, 0)
		inputView.SetOrigin(0, 0)

		return nil
	})

}
func writeToView(viewName string, message string) {
	g.Update(func(g *gocui.Gui) error {
		updatingView, err := g.View(viewName)
		if err != nil {
			return err
		}
		for _, c := range message {
			updatingView.EditWrite(c)
		}

		return nil
	})
}

// this removes formatting
func printError(message string) {
	printErrorF(message)
}
func printErrorF(message string, parts ...StyledString) {
	printToView("Feed", config.Colors.Feed.Error.sprintf(removeFormatting(message), parts...).string())
}

// this removes formatting
func printInfo(message string) {
	printInfoF(message)
}

// this removes formatting
func printInfoF(message string, parts ...StyledString) {
	printToView("Feed", config.Colors.Feed.Basic.sprintf(removeFormatting(message), parts...).string())
}
func printToView(viewName string, message string) {
	g.Update(func(g *gocui.Gui) error {
		updatingView, err := g.View(viewName)
		if err != nil {
			return err
		} else {
			if config.Basics.UnicodeEmojis {
				message = emojiUnicodeConvert(message)
			}
			fmt.Fprintf(updatingView, "%s\n", message)
		}
		return nil
	})
}

// End gocui helper funcs

// Update/Populate views automatically
func updateChatWindow() {

	runOpts := keybase.RunOptions{
		Dev: dev,
	}
	k.Run(func(api keybase.ChatAPI) {
		handleMessage(api)
	},
		runOpts)

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
			printError(fmt.Sprintf("%+v", err))
			return
		}
		go populateChat()
		go generateChannelTabCompletionSlice()
		return
	}
	var printMe []string
	var actuallyPrintMe string
	if len(api.Result.Messages) > 0 {
		lastMessage.ID = api.Result.Messages[0].Msg.ID
	}
	for _, message := range api.Result.Messages {
		if message.Msg.Content.Type == "text" || message.Msg.Content.Type == "attachment" {
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
func populateList() {
	_, maxY := g.Size()
	if testVar, err := k.ChatList(); err != nil {
		log.Printf("%+v", err)
	} else {
		clearView("List")
		var textBase = config.Colors.Channels.Basic.stylize("")
		var recentPMs = textBase.append(config.Colors.Channels.Header.stylize("---[PMs]---\n"))
		var recentPMsCount = 0
		var recentChannels = textBase.append(config.Colors.Channels.Header.stylize("---[Teams]---\n"))
		var recentChannelsCount = 0
		for _, s := range testVar.Result.Conversations {
			channels = append(channels, s.Channel)
			if s.Channel.MembersType == keybase.TEAM {
				recentChannelsCount++
				if recentChannelsCount <= ((maxY - 2) / 3) {
					channel := fmt.Sprintf("%s\n\t#%s\n", s.Channel.Name, s.Channel.TopicName)
					if s.Unread {
						recentChannels = recentChannels.append(config.Colors.Channels.Unread.stylize("*" + channel))
					} else {
						recentChannels = recentChannels.appendString(channel)
					}
				}
			} else {
				recentPMsCount++
				if recentPMsCount <= ((maxY - 2) / 3) {
					pmName := fmt.Sprintf("%s\n", cleanChannelName(s.Channel.Name))
					if s.Unread {
						recentPMs = recentPMs.append(config.Colors.Channels.Unread.stylize("*" + pmName))
					} else {
						recentPMs = recentPMs.appendString(pmName)
					}
				}
			}
		}
		time.Sleep(1 * time.Millisecond)
		printToView("List", fmt.Sprintf("%s%s", recentPMs.string(), recentChannels.string()))
		generateRecentTabCompletionSlice()
	}
}

// End update/populate views automatically

// Formatting
func formatMessageBody(body string) StyledString {
	output := config.Colors.Message.Body.stylize(body)

	output = colorReplaceMentionMe(output)
	output = output.colorRegex(`_[^_]*_`, config.Colors.Message.Body.withItalic())
	output = output.colorRegex(`~[^~]*~`, config.Colors.Message.Body.withStrikethrough())
	output = output.colorRegex(`@[\w_]*(\.[\w_]+)*`, config.Colors.Message.LinkKeybase)
	// TODO change how bold, italic etc works, so it uses boldOn boldOff ([1m and [22m)
	output = output.colorRegex(`\*[^\*]*\*`, config.Colors.Message.Body.withBold())
	output = output.replaceString("```", "<code>")
	// TODO make background color cover whole line
	output = output.colorRegex("<code>(.*\n)*<code>", config.Colors.Message.Code)
	output = output.colorRegex("`[^`]*`", config.Colors.Message.Code)
	// mention URL
	output = output.colorRegex(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*))`, config.Colors.Message.LinkURL)
	return output
}

// TODO use this more
func formatChannel(ch keybase.Channel) StyledString {
	return config.Colors.Message.LinkKeybase.stylize(fmt.Sprintf("@%s#%s", ch.Name, ch.TopicName))
}

func colorReplaceMentionMe(msg StyledString) StyledString {
	return msg.colorRegex(`(@?\b`+k.Username+`\b)`, config.Colors.Message.Mention)
}
func colorUsername(username string) StyledString {
	var color = config.Colors.Message.SenderDefault
	if username == k.Username {
		color = config.Colors.Message.Mention
	}
	return color.stylize(username)
}

func cleanChannelName(c string) string {
	newChannelName := strings.Replace(c, fmt.Sprintf("%s,", k.Username), "", 1)
	return strings.Replace(newChannelName, fmt.Sprintf(",%s", k.Username), "", 1)
}

func formatMessage(api keybase.ChatAPI, formatString string) string {
	ret := config.Colors.Message.Header.stylize("")
	msgType := api.Msg.Content.Type
	switch msgType {
	case "text", "attachment":
		ret = config.Colors.Message.Header.stylize(formatString)
		tm := time.Unix(int64(api.Msg.SentAt), 0)
		var msg = formatMessageBody(api.Msg.Content.Text.Body)
		if msgType == "attachment" {
			msg = config.Colors.Message.Body.stylize("$TITLE\n$FILE")
			attachment := api.Msg.Content.Attachment
			msg = msg.replaceString("$TITLE", attachment.Object.Title)
			msg = msg.replace("$FILE", config.Colors.Message.Attachment.stylize(fmt.Sprintf("[Attachment: %s]", attachment.Object.Filename)))
		}

		user := colorUsername(api.Msg.Sender.Username)
		device := config.Colors.Message.SenderDevice.stylize(api.Msg.Sender.DeviceName)
		msgID := config.Colors.Message.ID.stylize(fmt.Sprintf("%d", api.Msg.ID))
		date := config.Colors.Message.Time.stylize(tm.Format(config.Formatting.DateFormat))
		msgTime := config.Colors.Message.Time.stylize(tm.Format(config.Formatting.TimeFormat))

		channelName := config.Colors.Message.ID.stylize(fmt.Sprintf("@%s#%s", api.Msg.Channel.Name, api.Msg.Channel.TopicName))
		ret = ret.replace("$MSG", msg)
		ret = ret.replace("$USER", user)
		ret = ret.replace("$DEVICE", device)
		ret = ret.replace("$ID", msgID)
		ret = ret.replace("$TIME", msgTime)
		ret = ret.replace("$DATE", date)
		ret = ret.replace("$TEAM", channelName)
	}
	return ret.string()
}
func formatOutput(api keybase.ChatAPI) string {
	format := config.Formatting.OutputFormat
	if stream {
		format = config.Formatting.OutputStreamFormat
	}
	return formatMessage(api, format)
}

// End formatting

// Input handling
func handleMessage(api keybase.ChatAPI) {
	if _, ok := typeCommands[api.Msg.Content.Type]; ok {
		if api.Msg.Channel.MembersType == channel.MembersType && cleanChannelName(api.Msg.Channel.Name) == channel.Name {
			if channel.MembersType == keybase.TEAM && channel.TopicName != api.Msg.Channel.TopicName {
			} else {
				go typeCommands[api.Msg.Content.Type].Exec(api)
			}
		}
	}
	if api.Msg.Content.Type == "text" || api.Msg.Content.Type == "attachment" {
		go populateList()
		msgSender := api.Msg.Sender.Username
		if !stream {
			if msgSender != k.Username {
				if api.Msg.Channel.MembersType == keybase.TEAM {
					topicName := api.Msg.Channel.TopicName
					for _, m := range api.Msg.Content.Text.UserMentions {
						if m.Text == k.Username {
							// We are in a team
							if topicName != channel.TopicName {
								printInfo(formatMessage(api, config.Formatting.OutputMentionFormat))
								fmt.Print("\a")
							}

							break
						}
					}
				} else {
					if msgSender != channel.Name {
						printInfo(formatMessage(api, config.Formatting.PMFormat))
						fmt.Print("\a")
					}

				}
			}
			if api.Msg.Channel.MembersType == channel.MembersType && cleanChannelName(api.Msg.Channel.Name) == channel.Name {
				if channel.MembersType == keybase.USER || channel.MembersType == keybase.TEAM && channel.TopicName == api.Msg.Channel.TopicName {
					printToView("Chat", formatOutput(api))
					chat := k.NewChat(channel)
					lastMessage.ID = api.Msg.ID
					chat.Read(api.Msg.ID)
				}
			}
		} else {
			if api.Msg.Channel.MembersType == keybase.TEAM {
				printToView("Chat", formatOutput(api))
			} else {
				printToView("Chat", formatMessage(api, config.Formatting.PMFormat))
			}
		}
	} else {
		//TODO: For edit/delete run this
		if api.Msg.Channel.MembersType == channel.MembersType && cleanChannelName(api.Msg.Channel.Name) == channel.Name {
			go populateChat()
		}
	}
}
func getInputString(viewName string) (string, error) {
	inputView, err := g.View(viewName)
	if err != nil {
		return "", err
	}
	retString := inputView.Buffer()
	retString = strings.Replace(retString, "\n", "", 800)
	return retString, err
}
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
func handleInput(viewName string) error {
	clearView(viewName)
	inputString, _ := getInputString(viewName)
	if inputString == "" {
		return nil
	}
	if strings.HasPrefix(inputString, config.Basics.CmdPrefix) {
		cmd := deleteEmpty(strings.Split(inputString[len(config.Basics.CmdPrefix):], " "))
		if len(cmd) < 1 {
			return nil
		}
		if c, ok := commands[cmd[0]]; ok {
			c.Exec(cmd)
			return nil
		} else if cmd[0] == "q" || cmd[0] == "quit" {
			return gocui.ErrQuit
		} else {
			printError(fmt.Sprintf("Command '%s' not recognized", cmd[0]))
			return nil
		}
	}
	if inputString[:1] == "+" || inputString[:1] == "-" {
		cmd := strings.Split(inputString, " ")
		cmd[0] = inputString[:1]
		RunCommand(cmd...)
	} else {
		inputString = resolveRootEmojis(inputString)
		go sendChat(inputString)
	}
	// restore any tab completion view titles on input commit
	if newViewTitle := getViewTitle(viewName); newViewTitle != "" {
		setViewTitle(viewName, newViewTitle)
	}

	go populateList()
	return nil
}
func sendChat(message string) {
	autoScrollView("Chat")
	chat := k.NewChat(channel)
	_, err := chat.Send(message)
	if err != nil {
		printError(fmt.Sprintf("There was an error %+v", err))
	}
}

// End input handling

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// RegisterTypeCommand registers a command to be used within the client
func RegisterTypeCommand(c TypeCommand) error {
	var notAdded string
	for _, cmd := range c.Cmd {
		if _, ok := typeCommands[cmd]; !ok {
			typeCommands[cmd] = c
			continue
		}
		notAdded = fmt.Sprintf("%s, %s", notAdded, cmd)
	}
	if notAdded != "" {
		return fmt.Errorf("The following aliases were not added because they already exist: %s", notAdded)
	}
	return nil
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
