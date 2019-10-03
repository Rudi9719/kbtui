// +build !rm_basic_commands allcommands downloadcmd

package main

import (
	"fmt"
	"strconv"
)

func init() {
	command := Command{
		Cmd:         []string{"download", "d"},
		Description: "Download a file",
		Help:        "",
		Exec:        cmdDownloadFile,
	}

	RegisterCommand(command)
}

func cmdDownloadFile(cmd []string) {
	messageID, _ := strconv.Atoi(cmd[1])
	fileName := cmd[2]

	chat := k.NewChat(channel)
	_, err := chat.Download(messageID, fmt.Sprintf("%s/%s", downloadPath, fileName))
	if err != nil {
		printToView("Feed", fmt.Sprintf("There was an error downloading %s from %s", fileName, channel.Name))
	} else {
		printToView("Feed", fmt.Sprintf("Downloaded %s from %s", fileName, channel.Name))
	}
}
