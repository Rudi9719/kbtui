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
	if len(cmd) < 2 {
		printInfo(fmt.Sprintf("%s%s $filePath $fileName - Upload file from absolute path with optional name", config.Basics.CmdPrefix, cmd[0]))
		return
	}
	filePath := cmd[1]
	if !strings.HasPrefix(filePath, "/") {
		dir, err := os.Getwd()
		if err != nil {
			printError(fmt.Sprintf("There was an error determining path %+v", err))
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
	channelName := config.Colors.Message.LinkKeybase.stylize(channel.Name)
	fileNameStylized := config.Colors.Feed.File.stylize(filePath)
	if err != nil {
		printError(fmt.Sprintf("There was an error uploading %s to %s\n%+v", filePath, channel.Name, err))
	} else {
		printInfoF("Uploaded $TEXT to $TEXT", fileNameStylized, channelName)
	}
}
