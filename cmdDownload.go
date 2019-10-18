// +build !rm_basic_commands allcommands downloadcmd

package main

import (
	"fmt"
	"strconv"
)

func init() {
	command := Command{
		Cmd:         []string{"download", "d"},
		Description: "$messageId $fileName - Download a file to user's downloadpath",
		Help:        "",
		Exec:        cmdDownloadFile,
	}

	RegisterCommand(command)
}

func cmdDownloadFile(cmd []string) {

	if len(cmd) < 2 {
		printToView("Feed", fmt.Sprintf("%s%s $messageId $fileName - Download a file to user's downloadpath", cmdPrefix, cmd[0]))
		return
	}
	messageID, _ := strconv.Atoi(cmd[1])
	var fileName string
	if len(cmd) == 3 {
		fileName = cmd[2]
	} else {
		fileName = ""
	}

	chat := k.NewChat(channel)
	_, err := chat.Download(messageID, fmt.Sprintf("%s/%s", downloadPath, fileName))
	if err != nil {
		printToView("Feed", fmt.Sprintf("There was an error downloading %s from %s", fileName, channel.Name))
	} else {
		printToView("Feed", fmt.Sprintf("Downloaded %s from %s", fileName, channel.Name))
	}
}
