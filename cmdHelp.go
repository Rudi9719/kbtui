// +build !rm_basic_commands allcommands helpcmd

package main

import (
	"fmt"
	"sort"

	"github.com/jroimartin/gocui"
)

func init() {
	command := Command{
		Cmd:         []string{"help", "h"},
		Description: "Show information about avaailable commands",
		Help:        "",
		Exec:        cmdHelp,
	}

	RegisterCommand(command)
}

func cmdHelp(g *gocui.Gui, cmd []string) {
	var helpText string
	if len(cmd) == 1 {
		sort.Strings(baseCommands)
		for _, c := range baseCommands {
			helpText = fmt.Sprintf("%s%s%s\t\t%s\n", helpText, cmdPrefix, c, commands[c].Description)
		}
	}
	printToView("Chat", helpText)
}
