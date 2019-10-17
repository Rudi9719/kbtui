package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"regexp"
	"strings"
)

// Begin Colors
type color int

const (
	black color = iota
	red
	green
	yellow
	purple
	magenta
	cyan
	grey
	normal color = -1
)

func colorFromString(s string) color {
	s = strings.ToLower(s)
	switch s {
	case "black":
		return black
	case "red":
		return red
	case "green":
		return green
	case "yellow":
		return yellow
	case "purple":
		return purple
	case "magenta":
		return magenta
	case "cyan":
		return cyan
	case "grey":
		return grey
	case "normal":
		return normal
	default:
		printError(fmt.Sprintf("color `%s` cannot be parsed.", s))
	}
	return normal
}

// Style struct for specializing the style/color of a stylize
type Style struct {
	foregroundColor color
	backgroundColor color
	bold            bool
	italic          bool // Currently not supported by the UI library
	underline       bool
	strikethrough   bool // Currently not supported by the UI library
	inverse         bool
}

var basicStyle = Style{normal, normal, false, false, false, false, false}

func styleFromConfig(config *toml.Tree, key string) Style {
	key = "Colors." + key + "."
	style := basicStyle
	if config.Has(key + "foreground") {
		style = style.withForeground(colorFromString(config.Get(key + "foreground").(string)))
	}
	if config.Has(key + "background") {
		style = style.withForeground(colorFromString(config.Get(key + "background").(string)))
	}
	if config.GetDefault(key+"bold", false).(bool) {
		style = style.withBold()
	}
	if config.GetDefault(key+"italic", false).(bool) {
		style = style.withItalic()
	}
	if config.GetDefault(key+"underline", false).(bool) {
		style = style.withUnderline()
	}
	if config.GetDefault(key+"strikethrough", false).(bool) {
		style = style.withStrikethrough()
	}
	if config.GetDefault(key+"inverse", false).(bool) {
		style = style.withInverse()
	}

	return style
}

func (s Style) withForeground(f color) Style {
	s.foregroundColor = f
	return s
}
func (s Style) withBackground(f color) Style {
	s.backgroundColor = f
	return s
}
func (s Style) withBold() Style {
	s.bold = true
	return s
}
func (s Style) withInverse() Style {
	s.inverse = true
	return s
}
func (s Style) withItalic() Style {
	s.italic = true
	return s
}
func (s Style) withStrikethrough() Style {
	s.strikethrough = true
	return s
}
func (s Style) withUnderline() Style {
	s.underline = true
	return s
}

// TODO create both as `reset` (which it is now) as well as `append`
//  which essentially just adds on top. that is relevant in the case of
//  bold/italic etc - it should add style - not clear.
func (s Style) toANSI() string {
	if colorless {
		return ""
	}
	output := "\x1b[0m\x1b[0"
	if s.foregroundColor != normal {
		output += fmt.Sprintf(";%d", 30+s.foregroundColor)
	}
	if s.backgroundColor != normal {
		output += fmt.Sprintf(";%d", 40+s.backgroundColor)
	}
	if s.bold {
		output += ";1"
	}
	if s.italic {
		output += ";3"
	}
	if s.underline {
		output += ";4"
	}
	if s.inverse {
		output += ";7"
	}
	if s.strikethrough {
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
