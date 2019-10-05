// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Build kbtui with just the basic commands.
func Build() {
	sh.Run("go", "build")
}

// Build kbtui with the basic commands, and the ShowReactions "TypeCommand".
// The ShowReactions TypeCommand will print a message in the feed window when
// a reaction is received in the current conversation.
func BuildShowReactions() {
	sh.Run("go", "build", "-tags", "showreactionscmd")
}

// Build kbtui with the basec commands, and the AutoReact "TypeCommand".
// The AutoReact TypeCommand will automatically react to every message
// received in the current conversation. This gets pretty annoying, and
// is not recommended.
func BuildAutoReact() {
	sh.Run("go", "build", "-tags", "autoreactcmd")
}

// Build kbtui with all Commands and TypeCommands enabled.
func BuildAllCommands() {
	sh.Run("go", "build", "-tags", "type_commands")
}
