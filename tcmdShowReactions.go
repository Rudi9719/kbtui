// +build !rm_basic_commands type_commands showreactionscmd

package main

import (
	"fmt"

	"samhofi.us/x/keybase"
)

func init() {
	command := TypeCommand{
		Cmd:         []string{"reaction"},
		Name:        "ShowReactions",
		Description: "Prints a message in the feed any time a reaction is received",
		Exec:        tcmdShowReactions,
	}

	RegisterTypeCommand(command)
}

func tcmdShowReactions(m keybase.ChatAPI) {
	where := ""
	team := false
	if m.Msg.Channel.MembersType == keybase.TEAM {
		team = true
		where = fmt.Sprintf("in @%s#%s", m.Msg.Channel.Name, m.Msg.Channel.TopicName)
	} else {
		where = fmt.Sprintf("in a PM")
	}
	printToView("Feed", fmt.Sprintf("%s reacted to %d with %s %s", m.Msg.Sender.Username, m.Msg.Content.Reaction.M, m.Msg.Content.Reaction.B, where))
	if channel.Name == m.Msg.Channel.Name {
		if team {
			if channel.TopicName == m.Msg.Channel.TopicName {
				clearView("Chat")
				go populateChat()
			}
			
		} else {
			clearView("Chat")
			go populateChat()
		}

	}

}
