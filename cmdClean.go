// +build !rm_basic_commands allcommands cleancmd

package main

func init() {
	command := Command{
		Cmd:         []string{"clean", "c"},
		Description: "- Clean, or redraw chat view",
		Help:        "",
		Exec:        cmdClean,
	}

	RegisterCommand(command)
}

func cmdClean(cmd []string) {
	clearView("Chat")
	clearView("List")
	go populateChat()
	go populateList()
}
