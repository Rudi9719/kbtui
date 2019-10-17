// +build !rm_basic_commands allcommands uploadcmd

package main

import (
	"fmt"
	"os"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"upload", "u"},
		Description: "$filePath $fileName - Upload file from absolute path with optional name",
		Help:        "",
		Exec:        cmdUploadFile,
	}

	RegisterCommand(command)
}

func cmdUploadFile(cmd []string) {
	filePath := cmd[1]
	if !strings.HasPrefix(filePath, "/") {
		dir, err := os.Getwd()
		if err != nil {
			printToView("Feed", fmt.Sprintf("There was an error determining path %+v", err))
		}
		filePath = fmt.Sprintf("%s/%s", dir, filePath)
	}
	var fileName string
	if len(cmd) == 3 {
		fileName = cmd[2]
	} else {
		fileName = ""
	}
	chat := k.NewChat(channel)
	_, err := chat.Upload(fileName, filePath)
	if err != nil {
		printToView("Feed", fmt.Sprintf("There was an error uploading %s to %s\n%+v", filePath, channel.Name, err))
	} else {
		printToView("Feed", fmt.Sprintf("Uploaded %s to %s", filePath, channel.Name))
	}
}
