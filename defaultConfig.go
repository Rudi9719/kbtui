package main

var defaultConfig = `
[basics]
downloadPath = "/tmp/"
colorless = false
# The prefix before evaluating a command
cmdPrefix = "/"

[formatting]
# BASH-like PS1 variable equivalent
outputFormat = "┌──[$USER@$DEVICE] [$ID] [$DATE - $TIME]\n└╼ $MSG"
outputStreamFormat = "┌──[$USER@$DEVICE] [$ID] [$DATE - $TIME]\n└╼ $MSG"
outputMentionFormat = "┌──[$USER@$DEVICE] [$ID] [$DATE - $TIME]\n└╼ $MSG"
pmFormat = "PM from $USER@$DEVICE: $MSG"

# 02 = Day, Jan = Month, 06 = Year
dateFormat = "02Jan06"

# 15 = hours, 04 = minutes, 05 = seconds
timeFormat = "15:04"


[colors]
	 [colors.channels]
		  [colors.channels.basic]
		  foreground = "normal"
		  [colors.channels.header]
		  foreground = "magenta"
		  bold = true
		  [colors.channels.unread]
		  foreground = "green"
		  italic = true

	 [colors.message]
		  [colors.message.body]
		  foreground = "normal"
		  [colors.message.header]
		  foreground = "grey"
		  [colors.message.mention]
		  foreground = "green"
		  italic = true
		  bold = true
		  [colors.message.id]
		  foreground = "yellow"
		  [colors.message.time]
		  foreground = "magenta"
		  [colors.message.sender_default]
		  foreground = "cyan"
		  bold = true
		  [colors.message.sender_device]
		  foreground = "cyan"
		  [colors.message.attachment]
		  foreground = "red"
		  [colors.message.link_url]
		  foreground = "yellow"
		  [colors.message.link_keybase]
		  foreground = "yellow"
		  [colors.message.reaction]
		  foreground = "magenta"
		  bold = true
		  [colors.message.code]
		  foreground = "cyan"
		  background = "grey"
	 [colors.feed]
		  [colors.feed.basic]
		  foreground = "grey"
		  [colors.feed.error]
		  foreground = "red"
`
