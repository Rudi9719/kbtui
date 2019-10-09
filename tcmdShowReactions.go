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
	printToView("Feed", fmt.Sprintf("%s reacted to %d with %s", m.Msg.Sender.Username, m.Msg.Content.Reaction.M, m.Msg.Content.Reaction.B))
}
