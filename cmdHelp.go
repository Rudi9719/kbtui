// +build !rm_basic_commands allcommands helpcmd

package main

import (
	"fmt"
	"sort"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"help", "h"},
		Description: "Show information about available commands",
		Help:        "",
		Exec:        cmdHelp,
	}

	RegisterCommand(command)
}

func cmdHelp(cmd []string) {
	var helpText string
	var tCommands []string
	if len(cmd) == 1 {
		sort.Strings(baseCommands)
		for _, c := range baseCommands {
			helpText = fmt.Sprintf("%s%s%s\t\t%s\n", helpText, cmdPrefix, c, commands[c].Description)
		}
		if len(typeCommands) > 0 {
			for c, _ := range typeCommands {
				tCommands = append(tCommands, typeCommands[c].Name)
			}
			sort.Strings(tCommands)
			helpText = fmt.Sprintf("%s\nThe following Type Commands are currently loaded: %s", helpText, strings.Join(tCommands, ", "))
		}
	}
	printToView("Chat", helpText)
}
