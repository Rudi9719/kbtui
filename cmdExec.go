// +build !rm_basic_commands allcommands execcmd

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"exec", "ex"},
		Description: "$keybase args - executes keybase $args and returns the output",
		Help:        "",
		Exec:        cmdExec,
	}
	RegisterCommand(command)
}

func cmdExec(cmd []string) {
	l := len(cmd)
	switch {
	case l >= 2:
		if cmd[1] == "keybase" {
			// if the user types /exec keybase wallet list
			// only send ["wallet", "list"]
			runKeybaseExec(cmd[2:])
		} else {
			// send everything except the command
			runKeybaseExec(cmd[1:])
		}
	case l == 1:
		fallthrough
	default:
		printExecHelp()
	}
}

func runKeybaseExec(args []string) {
	cmd := exec.Command("keybase", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		printToView("Feed", fmt.Sprintf("Exec error: %+v", err))
	} else {
		channel.Name = ""
		// unjoin the chat
		clearView("Chat")
		setViewTitle("Input", fmt.Sprintf(" /exec %s ", strings.Join(args, " ")))
		printToView("Chat", fmt.Sprintf("%s", output))
	}
}

func printExecHelp() {
	printInfo(fmt.Sprintf("To execute a keybase command use %sexec <keybase args>", config.Basics.CmdPrefix))
}
