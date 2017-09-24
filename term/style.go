package term

import termbox "github.com/nsf/termbox-go"

// Style represents a complete text style, including both foreground and
// background color. It is encoded as a 32-bit integer.
//
// The coding is (MSB): <16b background><16b foreground>.
//
// Styles are saved along with the foreground color.
type Style uint32

const (
	AttrBold      = Style(termbox.AttrBold)
	AttrReverse   = Style(termbox.AttrReverse)
	AttrUnderline = Style(termbox.AttrUnderline)
)

const (
	ColorDefault Style = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

var ColorNames = map[string]Style{
	"default":       ColorDefault,
	"black":         ColorBlack,
	"red":           ColorRed,
	"green":         ColorGreen,
	"yellow":        ColorYellow,
	"blue":          ColorBlue,
	"magenta":       ColorMagenta,
	"cyan":          ColorCyan,
	"white":         ColorWhite,
	"brightBlack":   ColorBrightBlack,
	"brightRed":     ColorBrightRed,
	"brightGreen":   ColorBrightGreen,
	"brightYellow":  ColorBrightYellow,
	"brightBlue":    ColorBrightBlue,
	"brightMagenta": ColorBrightMagenta,
	"brightCyan":    ColorBrightCyan,
	"brightWhite":   ColorBrightWhite,

	"gray": ColorBrightBlack,
	"grey": ColorBrightBlack,
}

// GetColor returns a style based on a text string.
func GetColor(name string) Style {
	return ColorNames[name]
}

// Foreground returns a new style based on s, with the foreground color set
// as requested.
func (s Style) Foreground(c Style) Style {
	return (s & 0xffff0000) | (c & 0x0000ffff)
}

// Background returns a new style based on s, with the background color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Background(c Style) Style {
	return (c << 16) | (s & 0x0000ffff)
}

func (s Style) setAttrs(attrs Style, on bool) Style {
	if on {
		return s | attrs
	}
	return s &^ attrs
}

// Bold returns a new style based on s, with the bold attribute set
// as requested.
func (s Style) Bold(on bool) Style {
	return s.setAttrs(AttrBold, on)
}

// Reverse returns a new style based on s, with the reverse attribute set
// as requested.  (Reverse usually changes the foreground and background
// colors.)
func (s Style) Reverse(on bool) Style {
	return s.setAttrs(AttrReverse, on)
}

// Underline returns a new style based on s, with the underline attribute set
// as requested.
func (s Style) Underline(on bool) Style {
	return s.setAttrs(AttrUnderline, on)
}

// Attr returns foreground and background attributes usable for termbox.
func (s Style) Attr() (fg, bg termbox.Attribute) {
	fg = termbox.Attribute(s & 0xffff)
	bg = termbox.Attribute((s >> 16) & 0xffff)
	return
}
