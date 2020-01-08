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
	filepath   string     `toml:"-"` // filepath is not stored in the config file, but is written to the Config struct so it's known where the config was loaded from
	Basics     Basics     `toml:"basics"`
	Formatting Formatting `toml:"formatting"`
	Colors     Colors     `toml:"colors"`
}

// Basics holds the 'basics' section of the config file
type Basics struct {
	DownloadPath  string `toml:"download_path"`
	Colorless     bool   `toml:"colorless"`
	CmdPrefix     string `toml:"cmd_prefix"`
	UnicodeEmojis bool   `toml:"unicode_emojis"`
}

// Formatting holds the 'formatting' section of the config file
type Formatting struct {
	OutputFormat           string `toml:"output_format"`
	OutputStreamFormat     string `toml:"output_stream_format"`
	OutputMentionFormat    string `toml:"output_mention_format"`
	PMFormat               string `toml:"pm_format"`
	DateFormat             string `toml:"date_format"`
	TimeFormat             string `toml:"time_format"`
	IconFollowingUser      string `toml:"icon_following_user"`
	IconIndirectFollowUser string `toml:"icon_indirect_following_user"`
}

// Colors holds the 'colors' section of the config file
type Colors struct {
	Channels Channels `toml:"channels"`
	Message  Message  `toml:"message"`
	Feed     Feed     `toml:"feed"`
}

// Style holds basic style information
type Style struct {
	Foreground    string `toml:"foreground"`
	Background    string `toml:"background"`
	Italic        bool   `toml:"italic"`
	Bold          bool   `toml:"bold"`
	Underline     bool   `toml:"underline"`
	Strikethrough bool   `toml:"strikethrough"`
	Inverse       bool   `toml:"inverse"`
}

// Channels holds the style information for various elements of a channel
type Channels struct {
	Basic  Style `toml:"basic"`
	Header Style `toml:"header"`
	Unread Style `toml:"unread"`
}

// Message holds the style information for various elements of a message
type Message struct {
	Body          Style `toml:"body"`
	Header        Style `toml:"header"`
	Mention       Style `toml:"mention"`
	ID            Style `toml:"id"`
	Tags          Style `toml:"tags"`
	Time          Style `toml:"time"`
	SenderDefault Style `toml:"sender_default"`
	SenderDevice  Style `toml:"sender_device"`
	SenderTags    Style `toml:"sender_tags"`
	Attachment    Style `toml:"attachment"`
	LinkURL       Style `toml:"link_url"`
	LinkKeybase   Style `toml:"link_keybase"`
	Reaction      Style `toml:"reaction"`
	Quote         Style `toml:"quote"`
	Code          Style `toml:"code"`
}

// Feed holds the style information for various elements of the feed window
type Feed struct {
	Basic   Style `toml:"basic"`
	Error   Style `toml:"error"`
	File    Style `toml:"file"`
	Success Style `toml:"success"`
}
