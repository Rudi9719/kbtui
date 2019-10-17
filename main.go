package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
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
	g           *gocui.Gui
)

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
		fmt.Fprintln(feedView, "Feed Window - If you are mentioned or receive a PM it will show here")
	}
	if chatView, err2 := g.SetView("Chat", maxX/2-maxX/3, maxY/5+1, maxX-1, maxY-5, 0); err2 != nil {
		if !gocui.IsUnknownView(err2) {
			return err2
		}
		chatView.Autoscroll = true
		chatView.Wrap = true
		fmt.Fprintf(chatView, "Welcome %s!\n\nYour chats will appear here.\nSupported commands are as follows:\n\n", k.Username)
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
		inputView.Title = fmt.Sprintf(" Not in a chat - write `%sj` to join", cmdPrefix)
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
func initKeybindings() error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			input, err := getInputString("Input")
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
			return handleInput("Input")
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Input", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return handleTab()
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
	return nil
}

// End gocui basic setup

// Gocui helper funcs
func setViewTitle(viewName string, title string) {
	g.Update(func(g *gocui.Gui) error {
		updatingView, err := g.View(viewName)
		if err != nil {
			return err
		} else {
			updatingView.Title = title
		}
		return nil
	})
}
func getViewTitle(viewName string) string {
	view, err := g.View(viewName)
	if err != nil {
		// in case there is active tab completion, filter that to just the view title and not the completion options.
		writeToView("Feed", fmt.Sprintf("Error getting view title: %s", err))
		return ""
	} else {
		return strings.Split(view.Title, "||")[0]
	}
}
func popupView(viewName string) {
	_, err := g.SetCurrentView(viewName)
	if err != nil {
		printToView("Feed", fmt.Sprintf("%+v", err))
	}
	_, err = g.SetViewOnTop(viewName)
	if err != nil {
		printToView("Feed", fmt.Sprintf("%+v", err))
	}
	g.Update(func(g *gocui.Gui) error {
		updatingView, err := g.View(viewName)
		if err != nil {
			return err
		} else {
			viewX, viewY := updatingView.Size()
			updatingView.MoveCursor(viewX, viewY, true)
		}
		return nil

	})
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
func writeToView(viewName string, message string) {
	g.Update(func(g *gocui.Gui) error {
		updatingView, err := g.View(viewName)
		if err != nil {
			return err
		} else {
			for _, c := range message {
				updatingView.EditWrite(c)
			}
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
			printToView("Feed", fmt.Sprintf("%+v", err))
			return
		} else {
			go populateChat()
			return
		}
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
		var recentPMs = fmt.Sprintf("%s---[PMs]---%s\n", channelsHeaderColor, channelsColor)
		var recentPMsCount = 0
		var recentChannels = fmt.Sprintf("%s---[Teams]---%s\n", channelsHeaderColor, channelsColor)
		var recentChannelsCount = 0
		for _, s := range testVar.Result.Conversations {
			channels = append(channels, s.Channel)
			if s.Channel.MembersType == keybase.TEAM {
				recentChannelsCount++
				if recentChannelsCount <= ((maxY - 2) / 3) {
					if s.Unread {
						recentChannels += fmt.Sprintf("%s*", color(0))
					}
					recentChannels += fmt.Sprintf("%s\n\t#%s\n%s", s.Channel.Name, s.Channel.TopicName, channelsColor)
				}
			} else {
				recentPMsCount++
				if recentPMsCount <= ((maxY - 2) / 3) {
					if s.Unread {
						recentChannels += fmt.Sprintf("%s*", color(0))
					}
					recentPMs += fmt.Sprintf("%s\n%s", cleanChannelName(s.Channel.Name), channelsColor)
				}
			}
		}
		time.Sleep(1 * time.Millisecond)
		printToView("List", fmt.Sprintf("%s%s%s%s", channelsColor, recentPMs, recentChannels, noColor))
	}
}

// End update/populate views automatically

// Formatting
func cleanChannelName(c string) string {
	newChannelName := strings.Replace(c, fmt.Sprintf("%s,", k.Username), "", 1)
	return strings.Replace(newChannelName, fmt.Sprintf(",%s", k.Username), "", 1)
}
func formatOutput(api keybase.ChatAPI) string {
	ret := ""
	msgType := api.Msg.Content.Type
	switch msgType {
	case "text", "attachment":
		var c = messageHeaderColor
		ret = colorText(outputFormat, c, noColor)
		tm := time.Unix(int64(api.Msg.SentAt), 0)
		var msg = api.Msg.Content.Text.Body
		// mention teams or users
		msg = colorRegex(msg, `(@\w*(\.\w+)*)`, messageLinkColor, messageBodyColor)
		// mention URL
		msg = colorRegex(msg, `(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*))`, messageLinkColor, messageBodyColor)
		msg = colorText(colorReplaceMentionMe(msg, messageBodyColor), messageBodyColor, c)
		if msgType == "attachment" {
			msg = fmt.Sprintf("%s\n%s", msg, colorText("[Attachment]", messageAttachmentColor, c))
		}

		user := colorUsername(api.Msg.Sender.Username, c)
		device := colorText(api.Msg.Sender.DeviceName, messageSenderDeviceColor, c)
		msgId := colorText(fmt.Sprintf("%d", api.Msg.ID), messageIdColor, c)
		ts := colorText(fmt.Sprintf("%s", tm.Format(timeFormat)), messageTimeColor, c)
		ret = strings.Replace(ret, "$MSG", msg, 1)
		ret = strings.Replace(ret, "$USER", user, 1)
		ret = strings.Replace(ret, "$DEVICE", device, 1)
		ret = strings.Replace(ret, "$ID", msgId, 1)
		ret = strings.Replace(ret, "$TIME", ts, 1)
		ret = strings.Replace(ret, "$DATE", fmt.Sprintf("%s", tm.Format(dateFormat)), 1)
		ret = strings.Replace(ret, "```", fmt.Sprintf("\n<code>\n"), -1)
	}
	return ret
}

// End formatting

// Tab completion
func getCurrentChannelMembership() []string {
	var rs []string
	if channel.Name != "" {
		t := k.NewTeam(channel.Name)
		if testVar, err := t.MemberList(); err != nil {
			return rs // then this isn't a team, its a PM or there was an error in the API call
		} else {
			for _, m := range testVar.Result.Members.Owners {
				rs = append(rs, fmt.Sprintf("%+v", m.Username))
			}
			for _, m := range testVar.Result.Members.Admins {
				rs = append(rs, fmt.Sprintf("%+v", m.Username))
			}
			for _, m := range testVar.Result.Members.Writers {
				rs = append(rs, fmt.Sprintf("%+v", m.Username))
			}
			for _, m := range testVar.Result.Members.Readers {
				rs = append(rs, fmt.Sprintf("%+v", m.Username))
			}
		}
	}
	return rs
}
func filterStringSlice(ss []string, fv string) []string {
	var rs []string
	for _, s := range ss {
		if strings.HasPrefix(s, fv) {
			rs = append(rs, s)
		}
	}
	return rs
}
func longestCommonPrefix(ss []string) string {
	// cover the case where the slice has no or one members
	switch len(ss) {
	case 0:
		return ""
	case 1:
		return ss[0]
	}
	// all strings are compared by bytes here forward (TBD unicode normalization?)
	// establish min, max lenth members of the slice by iterating over the members
	min, max := ss[0], ss[0]
	for _, s := range ss[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}
	// then iterate over the characters from min to max, as soon as chars don't match return
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	// to cover the case where all members are equal, just return one
	return min
}
func stringRemainder(aStr, bStr string) string {
	var long, short string
	//figure out which string is longer
	switch {
	case len(aStr) < len(bStr):
		short = aStr
		long = bStr
	default:
		short = bStr
		long = aStr
	}
	// iterate over the strings using an external iterator so we don't lose the value
	i := 0
	for i < len(short) && i < len(long) {
		if short[i] != long[i] {
			// the strings aren't equal so don't return anything
			return ""
		}
		i++
	}
	// return whatever's left of the longer string
	return long[i:]
}
func appendIfNotInSlice(ss []string, s string) []string {
	for _, element := range ss {
		if element == s {
			return ss
		}
	}
	return append(ss, s)
}
func generateChannelTabCompletionSlice(inputWord string) []string {
	// create a slice to hold the values
	var firstSlice []string
	// iterate over all the conversation results
	for _, s := range channels {
		if s.MembersType == keybase.TEAM {
			// its a team so add the topic name as a possible tab completion
			firstSlice = appendIfNotInSlice(firstSlice, s.TopicName)
			firstSlice = appendIfNotInSlice(firstSlice, s.Name)
		} else {
			// its a user, so clean the name and append the users name as a possible tab completion
			firstSlice = appendIfNotInSlice(firstSlice, cleanChannelName(s.Name))
		}
	}
	// next fetch all members of the current channel and add them to the slice
	secondSlice := getCurrentChannelMembership()
	for _, m := range secondSlice {
		firstSlice = appendIfNotInSlice(firstSlice, m)
	}
	// now return the resultSlice which contains all that are prefixed with inputWord
	resultSlice := filterStringSlice(firstSlice, inputWord)
	return resultSlice
}
func generateEmojiTabCompletionSlice(inputWord string) []string {
	// use the emojiSlice from emojiList.go and filter it for the input word
	resultSlice := filterStringSlice(emojiSlice, inputWord)
	return resultSlice
}
func handleTab() error {
	inputString, err := getInputString("Input")
	if err != nil {
		return err
	} else {
		// if you successfully get an input string, grab the last word from the string
		ss := regexp.MustCompile(`[ #]`).Split(inputString, -1)
		s := ss[len(ss)-1]
		// create a variable in which to store the result
		var resultSlice []string
		// if the word starts with a : its an emoji lookup
		if strings.HasPrefix(s, ":") {
			resultSlice = generateEmojiTabCompletionSlice(s)
		} else {
			// now in case the word (s) is a mention @something, lets remove it to normalize
			if strings.HasPrefix(s, "@") {
				s = strings.Replace(s, "@", "", 1)
			}
			// now call get the list of all possible cantidates that have that as a prefix
			resultSlice = generateChannelTabCompletionSlice(s)
		}
		rLen := len(resultSlice)
		lcp := longestCommonPrefix(resultSlice)
		if lcp != "" {
			originalViewTitle := getViewTitle("Input")
			newViewTitle := ""
			if rLen >= 1 && originalViewTitle != "" {
				if rLen == 1 {
					newViewTitle = originalViewTitle
				} else if rLen <= 5 {
					newViewTitle = fmt.Sprintf("%s|| %s", originalViewTitle, strings.Join(resultSlice, " "))
				} else if rLen > 5 {
					newViewTitle = fmt.Sprintf("%s|| %s +%d more", originalViewTitle, strings.Join(resultSlice[:6], " "), rLen-5)
				}
				setViewTitle("Input", newViewTitle)
				remainder := stringRemainder(s, lcp)
				writeToView("Input", remainder)
			}
		}
	}
	return nil
}

// End tab completion

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
								fmt.Print("\a")
							}

							break
						}
					}
				} else {
					if msgSender != channel.Name {
						printToView("Feed", fmt.Sprintf("PM from @%s: %s", cleanChannelName(channelName), msgBody))
						fmt.Print("\a")
					}

				}
			}
			if api.Msg.Channel.MembersType == channel.MembersType && cleanChannelName(api.Msg.Channel.Name) == channel.Name {
				if channel.MembersType == keybase.TEAM && channel.TopicName == api.Msg.Channel.TopicName {
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
	if strings.HasPrefix(inputString, cmdPrefix) {
		cmd := deleteEmpty(strings.Split(inputString[len(cmdPrefix):], " "))
		if len(cmd) < 1 {
			return nil
		}
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
	if inputString[:1] == "+" || inputString[:1] == "-" {
		cmd := strings.Split(inputString, " ")
		RunCommand(cmd...)
	} else {
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
	chat := k.NewChat(channel)
	_, err := chat.Send(message)
	if err != nil {
		printToView("Feed", fmt.Sprintf("There was an error %+v", err))
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
