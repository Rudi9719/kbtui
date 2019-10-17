package main

// Path where Downloaded files will default to
var downloadPath = "/tmp/"

var colorless bool = false
var channelsColor = basicStyle
var channelUnreadColor = channelsColor.withForeground(green).withItalic()
var channelsHeaderColor = channelsColor.withForeground(magenta).withBold()

var mentionColor = basicStyle.withForeground(green)
var messageHeaderColor = basicStyle.withForeground(grey)
var messageIDColor = basicStyle.withForeground(yellow)
var messageTimeColor = basicStyle.withForeground(magenta)
var messageSenderDefaultColor = basicStyle.withForeground(cyan)
var messageSenderDeviceColor = messageSenderDefaultColor
var messageBodyColor = basicStyle
var messageAttachmentColor = basicStyle.withForeground(red)
var messageLinkURLColor = basicStyle.withForeground(yellow)
var messageLinkKeybaseColor = basicStyle.withForeground(yellow)
var messageReactionColor = basicStyle.withForeground(magenta)
var messageCodeColor = basicStyle.withBackground(grey).withForeground(cyan)

var feedColor = basicStyle.withForeground(grey)
var errorColor = basicStyle.withForeground(red)

// BASH-like PS1 variable equivalent
var outputFormat = "┌──[$USER@$DEVICE] [$ID] [$DATE - $TIME]\n└╼ $MSG"
var outputStreamFormat = "┌──[$TEAM] [$USER@$DEVICE] [$ID] [$DATE - $TIME]\n└╼ $MSG"
var mentionFormat = outputStreamFormat
var pmFormat = "PM from $USER@$DEVICE: $MSG"

// 02 = Day, Jan = Month, 06 = Year
var dateFormat = "02Jan06"

// 15 = hours, 04 = minutes, 05 = seconds
var timeFormat = "15:04"

// The prefix before evaluating a command
var cmdPrefix = "/"
