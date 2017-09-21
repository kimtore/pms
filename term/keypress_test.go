package term_test

import (
	"testing"

	"github.com/ambientsound/pms/term"
	termbox "github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

type keypressTest struct {
	input  termbox.Event
	output term.KeyPress
	format string
}

var keypressTests = []keypressTest{
	{
		termbox.Event{Ch: 'a'},
		term.KeyPress{0, 'a', 0},
		"a",
	},
	{
		termbox.Event{Key: termbox.KeySpace},
		term.KeyPress{termbox.KeySpace, ' ', 0},
		"<Space>",
	},
	{
		termbox.Event{Key: termbox.KeyCtrlA},
		term.KeyPress{termbox.KeyCtrlA, 'a', term.ModCtrl},
		"<Ctrl-A>",
	},
	{
		termbox.Event{Mod: termbox.ModAlt, Ch: 'e'},
		term.KeyPress{0, 'e', term.ModAlt},
		"<Alt-e>",
	},
	{
		termbox.Event{Mod: termbox.ModAlt, Ch: 'E'},
		term.KeyPress{0, 'E', term.ModAlt},
		"<Alt-E>",
	},
	{
		termbox.Event{Key: termbox.KeyF10},
		term.KeyPress{termbox.KeyF10, 0, 0},
		"<F10>",
	},
	{
		termbox.Event{Key: termbox.KeyF10, Mod: termbox.ModAlt},
		term.KeyPress{termbox.KeyF10, 0, term.ModAlt},
		"<Alt-F10>",
	},
	{
		termbox.Event{Key: termbox.KeyCtrlL, Mod: termbox.ModAlt},
		term.KeyPress{termbox.KeyCtrlL, 'l', term.ModCtrl | term.ModAlt},
		"<Alt-Ctrl-L>",
	},
}

func TestKeypress(t *testing.T) {
	for n, test := range keypressTests {
		t.Logf("### Test %d ###", n)

		output := term.ParseKey(test.input)
		assert.Equal(t, test.output, output)

		format := test.output.Name()
		assert.Equal(t, test.format, format)
	}
}
