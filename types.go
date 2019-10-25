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

// Config holds user-configurable values
type Config struct {
	filepath   string     `toml:"-"`
	Basics     Basics     `toml:"basics"`
	Formatting Formatting `toml:"formatting"`
	Colors     Colors     `toml:"colors"`
}

type Basics struct {
	DownloadPath string `toml:"downloadPath"`
	Colorless    bool   `toml:"colorless"`
	CmdPrefix    string `toml:"cmdPrefix"`
}

type Formatting struct {
	OutputFormat        string `toml:"outputFormat"`
	OutputStreamFormat  string `toml:"outputStreamFormat"`
	OutputMentionFormat string `toml:"outputMentionFormat"`
	PMFormat            string `toml:"pmFormat"`
	DateFormat          string `toml:"dateFormat"`
	TimeFormat          string `toml:"timeFormat"`
}

type Style struct {
	Foreground    string `toml:"foreground"`
	Background    string `toml:"background"`
	Italic        bool   `toml:"italic"`
	Bold          bool   `toml:"bold"`
	Underline     bool   `toml:"underline"`
	Strikethrough bool   `toml:"strikethrough"`
	Inverse       bool   `toml:"inverse"`
}

type Channels struct {
	Basic  Style `toml:"basic"`
	Header Style `toml:"header"`
	Unread Style `toml:"unread"`
}

type Message struct {
	Body          Style `toml:"body"`
	Header        Style `toml:"header"`
	Mention       Style `toml:"mention"`
	ID            Style `toml:"id"`
	Time          Style `toml:"time"`
	SenderDefault Style `toml:"sender_default"`
	SenderDevice  Style `toml:"sender_device"`
	Attachment    Style `toml:"attachment"`
	LinkURL       Style `toml:"link_url"`
	LinkKeybase   Style `toml:"link_keybase"`
	Reaction      Style `toml:"reaction"`
	Code          Style `toml:"code"`
}

type Feed struct {
	Basic Style `toml:"basic"`
	Error Style `toml:"error"`
}

type Colors struct {
	Channels Channels `toml:"channels"`
	Message  Message  `toml:"message"`
	Feed     Feed     `toml:"feed"`
}
