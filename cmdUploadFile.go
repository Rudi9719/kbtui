// +build !rm_basic_commands allcommands joincmd

package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func init() {
	command := Command{
		Cmd:         []string{"upload", "u"},
		Description: "Upload a file",
		Help:        "",
		Exec:        cmdUploadFile,
	}

	RegisterCommand(command)
}

func cmdUploadFile(g *gocui.Gui, cmd []string) {
	filePath := cmd[1]
	var fileName string
	if len(cmd) == 3 {
		fileName = cmd[2]
	} else {
		fileName = ""
	}
	chat := k.NewChat(channel)
	_, err := chat.Upload(fileName, filePath)
	if err != nil {
		printToView(g, "Feed", fmt.Sprintf("There was an error uploading %s to %s", filePath, channel.Name))
	} else {
		printToView(g, "Feed", fmt.Sprintf("Uploaded %s to %s", filePath, channel.Name))
	}
}
