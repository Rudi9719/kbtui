// +build mage

package main

import (
	"encoding/json"
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// emoji related constants
const emojiList = "https://raw.githubusercontent.com/CodeFreezr/emojo/master/db/v5/emoji-v5.json"
const emojiFileName = "emojiList.go"

// json parsing structure
type emoji struct {
	Num         int    `json:"No"`
	Emoji       string `json:"Emoji"`
	Category    string `json:"Category"`
	SubCategory string `json:"SubCategory"`
	Unicode     string `json:"Unicode"`
	Name        string `json:"Name"`
	Tags        string `json:"Tags"`
	Shortcode   string `json:"Shortcode"`
}

// This func downloaded and parses the emojis from online into a slice of all shortnames
// to be used as a lookup for tab completion for emojis
// this way the pull from GitHub only has to be done at build time.
func createEmojiSlice() ([]string, error) {
	result, err := http.Get(emojiList)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	emojiList, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	var emojis []emoji
	if err := json.Unmarshal(emojiList, &emojis); err != nil {
		return nil, err
	}

	var emojiSlice []string

	for _, emj := range emojis {
		if len(emj.Shortcode) == 0 || strings.Contains(emj.Shortcode, "_tone") {
			// dont add them
			continue
		}
		emojiSlice = append(emojiSlice, emj.Shortcode)
	}
	return emojiSlice, nil
}

// Build kbtui with emoji lookup support
func BuildEmoji() error {
	emojis, err := createEmojiSlice()
	if err != nil {
		return err
	}
	f, err := os.Create(emojiFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	fileContent := fmt.Sprintf("package main\n\nvar emojiSlice = %#v", emojis)
	_, err = f.WriteString(fileContent)
	if err != nil {
		return err
	}
	f.Sync()
	return nil
}

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

// Build kbtui with all commands and TypeCommands disabled.
func BuildAllCommands() {
	sh.Run("go", "build", "-tags", "allcommands")
}

// Build kbtui with all Commands and TypeCommands enabled.
func BuildAllCommandsT() {
	sh.Run("go", "build", "-tags", "type_commands,allcommands")
}

// Build kbtui with beta functionality
func BuildBeta() {
	mg.Deps(BuildEmoji)
	sh.Run("go", "build", "-tags", "allcommands,showreactionscmd,emojiList")
}
