// +ignore
// +build allcommands tagscmd

package main

func init() {
	command := Command{
		Cmd:         []string{"tags", "map"},
		Description: "$- Create map of users following users, to populate $TAGS",
		Help:        "",
		Exec:        cmdTags,
	}

	RegisterCommand(command)
}

func cmdTags(cmd []string) {
	go generateFollowersList()
}
