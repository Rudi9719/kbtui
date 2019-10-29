// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func getRemotePackages() error {
	var packages = []string{
		"samhofi.us/x/keybase",
		"github.com/awesome-gocui/gocui",
		"github.com/magefile/mage/mage",
		"github.com/magefile/mage/mg",
		"github.com/magefile/mage/sh",
		"github.com/pelletier/go-toml",
	}
	for _, p := range packages {
		if err := sh.Run("go", "get", "-u", p); err != nil {
			return err
		}
	}
	return nil
}

// proper error reporting and exit code
func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

// Build kbtui with just the basic commands.
func Build() {
	mg.Deps(getRemotePackages)
	if err := sh.Run("go", "build"); err != nil {
		defer func() {
			exit(err)
		}()
	}
}

// Build kbtui with the basic commands, and the ShowReactions "TypeCommand".
// The ShowReactions TypeCommand will print a message in the feed window when
// a reaction is received in the current conversation.
func BuildShowReactions() {
	mg.Deps(getRemotePackages)
	if err := sh.Run("go", "build", "-tags", "showreactionscmd"); err != nil {
		defer func() {
			exit(err)
		}()
	}
}

// Build kbtui with the basec commands, and the AutoReact "TypeCommand".
// The AutoReact TypeCommand will automatically react to every message
// received in the current conversation. This gets pretty annoying, and
// is not recommended.
func BuildAutoReact() {
	mg.Deps(getRemotePackages)
	if err := sh.Run("go", "build", "-tags", "autoreactcmd"); err != nil {
		defer func() {
			exit(err)
		}()
	}
}

// Build kbtui with all commands and TypeCommands disabled.
func BuildAllCommands() {
	mg.Deps(getRemotePackages)
	if err := sh.Run("go", "build", "-tags", "allcommands"); err != nil {
		defer func() {
			exit(err)
		}()
	}
}

// Build kbtui with all Commands and TypeCommands enabled.
func BuildAllCommandsT() {
	mg.Deps(getRemotePackages)
	if err := sh.Run("go", "build", "-tags", "type_commands allcommands"); err != nil {
		defer func() {
			exit(err)
		}()
	}
}

// Build kbtui with beta functionality
func BuildBeta() {
	mg.Deps(getRemotePackages)
	if err := sh.Run("go", "build", "-tags", "allcommands showreactionscmd tabcompletion execcmd"); err != nil {
		defer func() {
			exit(err)
		}()
	}
}
