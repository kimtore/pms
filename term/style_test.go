package term_test

import (
	"testing"

	"github.com/ambientsound/pms/term"
	termbox "github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

type styleTest struct {
	input  term.Style
	output term.Style
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
		term.Style(0).Foreground(term.ColorRed),
		0x00000002,
		termbox.ColorRed,
		termbox.ColorDefault,
	},
	{
		term.Style(0).Foreground(term.ColorRed).Background(term.ColorBlue),
		0x00050002,
		termbox.Attribute(term.ColorRed),
		termbox.Attribute(term.ColorBlue),
	},
	{
		term.Style(0).Foreground(term.ColorGreen).Bold(true),
		0x00000203,
		termbox.Attribute(term.ColorGreen) | termbox.AttrBold,
		termbox.ColorDefault,
	},
	{
		term.Style(0).Underline(true),
		0x00000400,
		termbox.AttrUnderline,
		termbox.ColorDefault,
	},
	{
		term.Style(0).Background(term.ColorGreen).Reverse(true),
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
