package main
// Path where Downloaded files will default to
var downloadPath = "/tmp/"

var colorless = false
var channelsColor = color(8)
var channelsHeaderColor = color(6)
var noColor = color(-1)
var mentionColor = color(3)
var messageHeaderColor = color(8)
var messageIdColor = color(7)
var messageTimeColor = color(6)
var messageSenderDefaultColor = color(8)
var messageSenderDeviceColor = color(8)
var messageBodyColor = noColor
var messageAttachmentColor = color(2)
var messageLinkColor = color(4)

// BASH-like PS1 variable equivalent (without colours)
var outputFormat = "┌──[$USER@$DEVICE] [$ID] [$DATE - $TIME]\n└╼ $MSG"

// 02 = Day, Jan = Month, 06 = Year
var dateFormat = "02Jan06"

// 15 = hours, 04 = minutes, 05 = seconds
var timeFormat = "15:04"

// The prefix before evaluating a command
var cmdPrefix = "/"
