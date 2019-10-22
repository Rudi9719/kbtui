// +build !rm_basic_commands allcommands streamcmd

package main

func init() {
	command := Command{
		Cmd:         []string{"stream", "s"},
		Description: "- Stream all incoming messages",
		Help:        "",
		Exec:        cmdStream,
	}

	RegisterCommand(command)
}

func cmdStream(cmd []string) {
	stream = true
	channel.Name = ""

	printToView("Feed", "You are now viewing the formatted stream")
	setViewTitle("Input", " Stream - Not in a chat /j to join ")
	clearView("Chat")
}
