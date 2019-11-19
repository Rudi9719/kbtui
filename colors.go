package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	black int = iota
	red
	green
	yellow
	purple
	magenta
	cyan
	grey
	normal int = -1
)

var colorMapString = map[string]int{
	"black":   black,
	"red":     red,
	"green":   green,
	"yellow":  yellow,
	"purple":  purple,
	"magenta": magenta,
	"cyan":    cyan,
	"grey":    grey,
	"normal":  normal,
}

var colorMapInt = map[int]string{
	black:   "black",
	red:     "red",
	green:   "green",
	yellow:  "yellow",
	purple:  "purple",
	magenta: "magenta",
	cyan:    "cyan",
	grey:    "grey",
	normal:  "normal",
}

func colorFromString(color string) int {
	var result int
	color = strings.ToLower(color)
	result, ok := colorMapString[color]
	if !ok {
		return normal
	}
	return result
}

func colorFromInt(color int) string {
	var result string
	result, ok := colorMapInt[color]
	if !ok {
		return "normal"
	}
	return result
}

var basicStyle = Style{
	Foreground:    colorMapInt[normal],
	Background:    colorMapInt[normal],
	Italic:        false,
	Bold:          false,
	Underline:     false,
	Strikethrough: false,
	Inverse:       false,
}

func (s Style) withForeground(color int) Style {
	s.Foreground = colorFromInt(color)
	return s
}
func (s Style) withBackground(color int) Style {
	s.Background = colorFromInt(color)
	return s
}

func (s Style) withBold() Style {
	s.Bold = true
	return s
}
func (s Style) withInverse() Style {
	s.Inverse = true
	return s
}
func (s Style) withItalic() Style {
	s.Italic = true
	return s
}
func (s Style) withStrikethrough() Style {
	s.Strikethrough = true
	return s
}
func (s Style) withUnderline() Style {
	s.Underline = true
	return s
}

// TODO create both as `reset` (which it is now) as well as `append`
//  which essentially just adds on top. that is relevant in the case of
//  bold/italic etc - it should add style - not clear.
func (s Style) toANSI() string {
	if config.Basics.Colorless {
		return ""
	}
	styleSlice := []string{"0"}

	if colorFromString(s.Foreground) != normal {
		styleSlice = append(styleSlice, fmt.Sprintf("%d", 30+colorFromString(s.Foreground)))
	}
	if colorFromString(s.Background) != normal {
		styleSlice = append(styleSlice, fmt.Sprintf("%d", 40+colorFromString(s.Background)))
	}
	if s.Bold {
		styleSlice = append(styleSlice, "1")
	}
	if s.Italic {
		styleSlice = append(styleSlice, "3")
	}
	if s.Underline {
		styleSlice = append(styleSlice, "4")
	}
	if s.Inverse {
		styleSlice = append(styleSlice, "7")
	}
	if s.Strikethrough {
		styleSlice = append(styleSlice, "9")
	}

	return "\x1b[" + strings.Join(styleSlice, ";") + "m"
}

// End Colors
// Begin StyledString

// StyledString is used to save a message with a style, which can then later be rendered to a string
type StyledString struct {
	message string
	style   Style
}

func (ss StyledString) withStyle(style Style) StyledString {
	return StyledString{ss.message, style}
}

// TODO change StyledString to have styles at start-end indexes.

// TODO handle all formatting types
func (s Style) sprintf(base string, parts ...StyledString) StyledString {
	text := s.stylize(removeFormatting(base))
	//TODO handle posibility to escape
	re := regexp.MustCompile(`\$TEXT`)
	for len(re.FindAllString(text.message, 1)) > 0 {
		part := parts[0]
		parts = parts[1:]
		text = text.replaceN("$TEXT", part, 1)
	}
	return text
}

func (s Style) stylize(msg string) StyledString {
	return StyledString{msg, s}
}
func (ss StyledString) stringFollowedByStyle(style Style) string {
	return ss.style.toANSI() + ss.message + style.toANSI()
}
func (ss StyledString) string() string {
	return ss.stringFollowedByStyle(basicStyle)
}

func (ss StyledString) replace(match string, value StyledString) StyledString {
	return ss.replaceN(match, value, -1)
}
func (ss StyledString) replaceN(match string, value StyledString, n int) StyledString {
	ss.message = strings.Replace(ss.message, match, value.stringFollowedByStyle(ss.style), n)
	return ss
}
func (ss StyledString) replaceString(match string, value string) StyledString {
	ss.message = strings.Replace(ss.message, match, value, -1)
	return ss
}

// Overrides current formatting
func (ss StyledString) colorRegex(match string, style Style) StyledString {
	return ss.regexReplaceFunc(match, func(subString string) string {
		return style.stylize(removeFormatting(subString)).stringFollowedByStyle(ss.style)
	})
}

// Replacer function takes the current match as input and should return how the match should be preseneted instead
func (ss StyledString) regexReplaceFunc(match string, replacer func(string) string) StyledString {
	re := regexp.MustCompile(match)
	locations := re.FindAllStringIndex(ss.message, -1)
	var newMessage string
	var prevIndex int
	for _, loc := range locations {
		newSubstring := replacer(ss.message[loc[0]:loc[1]])
		newMessage += ss.message[prevIndex:loc[0]]
		newMessage += newSubstring
		prevIndex = loc[1]
	}
	// Append any string after the final match
	newMessage += ss.message[prevIndex:len(ss.message)]
	ss.message = newMessage
	return ss
}

// Appends the other stylize at the end, but retains same style
func (ss StyledString) append(other StyledString) StyledString {
	ss.message = ss.message + other.stringFollowedByStyle(ss.style)
	return ss
}
func (ss StyledString) appendString(other string) StyledString {
	ss.message += other
	return ss
}

// Begin Formatting

func removeFormatting(s string) string {
	reFormatting := regexp.MustCompile(`(?m)\x1b\[(\d*;?)*m`)
	return reFormatting.ReplaceAllString(s, "")
}
