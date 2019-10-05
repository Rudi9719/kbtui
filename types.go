package main

import "samhofi.us/x/keybase"

// Command outlines a command
type Command struct {
	Cmd         []string       // Any aliases that trigger this command
	Description string         // A short description of the command
	Help        string         // The full help text explaining how to use the command
	Exec        func([]string) // A function that takes the command (arg[0]) and any arguments (arg[1:]) as input
}

// TypeCommand outlines a command that reacts on message type
type TypeCommand struct {
	Cmd         []string              // Message types that trigger this command
	Name        string                // The name of this command
	Description string                // A short description of the command
	Exec        func(keybase.ChatAPI) // A function that takes a raw chat message as input
}
