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
	output := "\x1b[0m\x1b[0"
	if colorFromString(s.Foreground) != normal {
		output += fmt.Sprintf(";%d", 30+colorFromString(s.Foreground))
	}
	if colorFromString(s.Background) != normal {
		output += fmt.Sprintf(";%d", 40+colorFromString(s.Background))
	}
	if s.Bold {
		output += ";1"
	}
	if s.Italic {
		output += ";3"
	}
	if s.Underline {
		output += ";4"
	}
	if s.Inverse {
		output += ";7"
	}
	if s.Strikethrough {
		output += ";9"
	}

	return output + "m"
}

// End Colors
// Begin StyledString

// StyledString is used to save a message with a style, which can then later be rendered to a string
type StyledString struct {
	message string
	style   Style
}

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
func (t StyledString) stringFollowedByStyle(style Style) string {
	return t.style.toANSI() + t.message + style.toANSI()
}
func (t StyledString) string() string {
	return t.stringFollowedByStyle(basicStyle)
}

func (t StyledString) replace(match string, value StyledString) StyledString {
	return t.replaceN(match, value, -1)
}
func (t StyledString) replaceN(match string, value StyledString, n int) StyledString {
	t.message = strings.Replace(t.message, match, value.stringFollowedByStyle(t.style), n)
	return t
}
func (t StyledString) replaceString(match string, value string) StyledString {
	t.message = strings.Replace(t.message, match, value, -1)
	return t
}
func (t StyledString) replaceRegex(match string, value StyledString) StyledString {
	var re = regexp.MustCompile("(" + match + ")")
	t.message = re.ReplaceAllString(t.message, value.stringFollowedByStyle(t.style))
	return t
}

// Overrides current formatting
func (t StyledString) colorRegex(match string, style Style) StyledString {
	re := regexp.MustCompile("(" + match + ")")
	subStrings := re.FindAllString(t.message, -1)
	for _, element := range subStrings {
		cleanSubstring := style.stylize(removeFormatting(element))
		t.message = strings.Replace(t.message, element, cleanSubstring.stringFollowedByStyle(t.style), -1)
	}
	return t
	// Old versionreturn t.replaceRegex(match, style.stylize(`$1`))
}

// Appends the other stylize at the end, but retains same style
func (t StyledString) append(other StyledString) StyledString {
	t.message = t.message + other.stringFollowedByStyle(t.style)
	return t
}
func (t StyledString) appendString(other string) StyledString {
	t.message += other
	return t
}

// Begin Formatting

func removeFormatting(s string) string {
	reFormatting := regexp.MustCompile(`(?m)\x1b\[(\d*;?)*m`)
	return reFormatting.ReplaceAllString(s, "")
}
