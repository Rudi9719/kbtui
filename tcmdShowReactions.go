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
	team := false
	user := colorUsername(m.Msg.Sender.Username)
	id := config.Colors.Message.ID.stylize(fmt.Sprintf("%d", m.Msg.Content.Reaction.M))
	reaction := config.Colors.Message.Reaction.stylize(m.Msg.Content.Reaction.B)
	where := config.Colors.Message.LinkKeybase.stylize("a PM")
	if m.Msg.Channel.MembersType == keybase.TEAM {
		team = true
		where = formatChannel(m.Msg.Channel)
	} else {
	}
	printInfoF("$TEXT reacted to [$TEXT] with $TEXT in $TEXT", user, id, reaction, where)
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
