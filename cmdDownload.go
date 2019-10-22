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
	messageID, err := strconv.Atoi(cmd[1])
	if err != nil {
		printToView("Feed", "There was an error converting your messageID to an int")
		return
	}
	chat := k.NewChat(channel)
	api, err := chat.ReadMessage(messageID)
	if err != nil {
		printToView("Feed", fmt.Sprintf("There was an error pulling message %d", messageID))
		return
	}
	if api.Result.Messages[0].Msg.Content.Type != "attachment" {
		printToView("Feed", "No attachment detected")
		return
	}
	var fileName string
	if len(cmd) == 3 {
		fileName = cmd[2]
	} else {
		fileName = api.Result.Messages[0].Msg.Content.Attachment.Object.Filename
	}

	_, err = chat.Download(messageID, fmt.Sprintf("%s/%s", downloadPath, fileName))
	if err != nil {
		printToView("Feed", fmt.Sprintf("There was an error downloading %s from %s", fileName, channel.Name))
	} else {
		printToView("Feed", fmt.Sprintf("Downloaded %s from %s", fileName, channel.Name))
	}
}
