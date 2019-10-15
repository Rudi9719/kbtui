package main

import (
	"fmt"
    "regexp"
)

// TODO maybe datastructure 
// BASH-like PS1 variable equivalent (without colours)
// TODO bold? cursive etc?
func color(c int) string {
	if colorless {
		return ""
	}
	if c < 0 {
		return "\033[0m"
	} else {
		return fmt.Sprintf("\033[0;%dm", 29+c)
	}
}
// TODO maybe make the text into some datastructure which remembers the color
func colorText(text string, color string, offColor string) string {
	return fmt.Sprintf("%s%s%s", color, text, offColor)
}

func colorUsername(username string, offColor string) string {
	var color = messageSenderDefaultColor
	if username == k.Username {
		color = mentionColor
	}
	return colorText(username, color, offColor)
}
func colorRegex(msg string, match string, color string, offColor string) string {
	var re = regexp.MustCompile(match)
	return re.ReplaceAllString(msg, colorText(`$1`, color, offColor))
}

func colorReplaceMentionMe(msg string, offColor string) string {
	//var coloredOwnName = colorText(k.Username, mentionColor, offColor)
	//return strings.Replace(msg, k.Username, coloredOwnName, -1)
	return colorRegex(msg, "(@?"+k.Username+")", mentionColor, offColor)
}
