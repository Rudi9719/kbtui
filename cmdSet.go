// +build !rm_basic_commands allcommands setcmd

package main

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml"
)

func init() {
	command := Command{
		Cmd:         []string{"set", "config"},
		Description: "Change various settings",
		Help:        "",
		Exec:        cmdSet,
	}

	RegisterCommand(command)
}
func printSetting(cmd []string) {
	switch cmd[1] {
	case "load":
		loadFromToml()
		printInfo("Loading config from toml")
	case "downloadPath":
		printInfo(fmt.Sprintf("Setting for %s -> %s", cmd[1], downloadPath))
	case "outputFormat":
		printInfo(fmt.Sprintf("Setting for %s -> %s", cmd[1], outputFormat))
	case "dateFormat":
		printInfo(fmt.Sprintf("Setting for %s -> %s", cmd[1], dateFormat))
	case "timeFormat":
		printInfo(fmt.Sprintf("Setting for %s -> %s", cmd[1], timeFormat))
	case "cmdPrefix":
		printInfo(fmt.Sprintf("Setting for %s -> %s", cmd[1], cmdPrefix))
	default:
		printError(fmt.Sprintf("Unknown config value %s", cmd[1]))
	}
}
func cmdSet(cmd []string) {
	if len(cmd) < 2 {
		printError("No config value specified")
		return
	}
	if len(cmd) < 3 {
		printSetting(cmd)
		return
	}
	switch cmd[1] {
	case "downloadPath":
		if len(cmd) != 3 {
			printError("Invalid download path.")
		}
		downloadPath = cmd[2]
	case "outputFormat":
		outputFormat = strings.Join(cmd[1:], " ")
	case "dateFormat":
		dateFormat = strings.Join(cmd[1:], " ")
	case "timeFormat":
		timeFormat = strings.Join(cmd[1:], " ")
	case "cmdPrefix":
		cmdPrefix = cmd[2]
	default:
		printError(fmt.Sprintf("Unknown config value %s", cmd[1]))
	}

}
func loadFromToml() {
	config, err := toml.LoadFile("kbtui.tml")
	if err != nil {
		printError(fmt.Sprintf("Could not read config file: %+v", err))
		return
	}
	colorless = config.GetDefault("Basics.colorless", false).(bool)
	if config.Has("Basics.colorless") {
		colorless = config.Get("Basics.colorless").(bool)
	}
	if config.Has("Basics.downloadPath") {
		downloadPath = config.Get("Basics.downloadPath").(string)
	}
	if config.Has("Basics.cmdPrefix") {
		cmdPrefix = config.Get("Basics.cmdPrefix").(string)
	}
	if config.Has("Formatting.outputFormat") {
		outputFormat = config.Get("Formatting.outputFormat").(string)
	}
	if config.Has("Formatting.dateFormat") {
		dateFormat = config.Get("Formatting.dateFormat").(string)
	}
	if config.Has("Formatting.timeFormat") {
		timeFormat = config.Get("Formatting.timeFormat").(string)
	}
	channelsColor = styleFromConfig(config, "channels.basic")

	channelsHeaderColor = styleFromConfig(config, "channels.header")
	channelUnreadColor = styleFromConfig(config, "channels.unread")

	mentionColor = styleFromConfig(config, "message.mention")
	messageHeaderColor = styleFromConfig(config, "message.header")
	messageIDColor = styleFromConfig(config, "message.id")
	messageTimeColor = styleFromConfig(config, "message.time")
	messageSenderDefaultColor = styleFromConfig(config, "message.sender_default")
	messageSenderDeviceColor = styleFromConfig(config, "message.sender_device")
	messageBodyColor = styleFromConfig(config, "message.body")
	messageAttachmentColor = styleFromConfig(config, "message.attachment")
	messageLinkURLColor = styleFromConfig(config, "message.link_url")
	messageLinkKeybaseColor = styleFromConfig(config, "message.link_keybase")
	messageReactionColor = styleFromConfig(config, "message.reaction")
	messageCodeColor = styleFromConfig(config, "message.code")

	feedColor = styleFromConfig(config, "feed.basic")
	errorColor = styleFromConfig(config, "feed.error")

	RunCommand("clean")
}
