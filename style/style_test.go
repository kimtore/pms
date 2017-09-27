package style_test

import (
	"testing"

	"github.com/ambientsound/pms/style"
	termbox "github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

type styleTest struct {
	input  style.Style
	output style.Style
	fg     termbox.Attribute
	bg     termbox.Attribute
}

var styleTests = []styleTest{
	{
		0,
		0x00000000,
		termbox.ColorDefault,
		termbox.ColorDefault,
	},
	{
		style.Style(0).Foreground(style.ColorRed),
		0x00000002,
		termbox.ColorRed,
		termbox.ColorDefault,
	},
	{
		style.Style(0).Foreground(style.ColorRed).Background(style.ColorBlue),
		0x00050002,
		termbox.Attribute(style.ColorRed),
		termbox.Attribute(style.ColorBlue),
	},
	{
		style.Style(0).Foreground(style.ColorGreen).Bold(true),
		0x00000203,
		termbox.Attribute(style.ColorGreen) | termbox.AttrBold,
		termbox.ColorDefault,
	},
	{
		style.Style(0).Underline(true),
		0x00000400,
		termbox.AttrUnderline,
		termbox.ColorDefault,
	},
	{
		style.Style(0).Background(style.ColorGreen).Reverse(true),
		0x00030800,
		termbox.AttrReverse,
		termbox.ColorGreen,
	},
}

func TestStyle(t *testing.T) {
	for n, test := range styleTests {
		t.Logf("### Test %d ###", n)
		fg, bg := test.input.Attr()
		assert.Equal(t, test.output, test.input)
		assert.Equal(t, test.fg, fg)
		assert.Equal(t, test.bg, bg)
	}
}
