package main

import (
	"github.com/jroimartin/gocui"
)

// Command outlines a command
type Command struct {
	Cmd         []string                   // Any aliases that trigger this command
	Description string                     // A short description of the command
	Help        string                     // The full help text explaining how to use the command
	Exec        func(*gocui.Gui, []string) // A function that takes the command (arg[0]) and any arguments (arg[1:]) as input
}
