package main

// Path where Downloaded files will default to
var downloadPath = "/tmp/"

// BASH-like PS1 variable equivalent (without colours)
var outputFormat = "┌──[$USER@$DEVICE] [$ID] [$DATE - $TIME]\n└╼ $MSG"

// 02 = Day, Jan = Month, 06 = Year
var dateFormat = "02Jan06"

// 15 = hours, 04 = minutes, 05 = seconds
var timeFormat = "15:04"

// The prefix before evaluating a command
var cmdPrefix = "/"
