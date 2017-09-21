package keysequence_test

import (
	"strings"
	"testing"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/keysequence"
	"github.com/ambientsound/pms/term"
	termbox "github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

var parserTests = []struct {
	input   string
	output  string
	success bool
	keyseq  keysequence.KeySequence
}{
	{
		"abc",
		"abc",
		true,
		keysequence.KeySequence{
			{0, 'a', 0},
			{0, 'b', 0},
			{0, 'c', 0},
		},
	},
	{
		"<C-c>x<S-m-A-f1>f",
		"<Ctrl-C>x<Shift-Alt-Meta-F1>f",
		true,
		keysequence.KeySequence{
			{termbox.KeyCtrlC, 'c', term.ModCtrl},
			{0, 'x', 0},
			{termbox.KeyF1, 0, term.ModShift | term.ModMeta | term.ModAlt},
			{0, 'f', 0},
		},
	},
	{
		"<a-x>",
		"<Alt-x>",
		true,
		keysequence.KeySequence{
			{0, 'x', term.ModAlt},
		},
	},
	{
		"<c-x>x",
		"<Ctrl-X>x",
		true,
		keysequence.KeySequence{
			{termbox.KeyCtrlX, 'x', term.ModCtrl},
			{0, 'x', 0},
		},
	},
	{
		"<Space>",
		"<Space>",
		true,
		keysequence.KeySequence{
			{termbox.KeySpace, ' ', 0},
		},
	},
	{
		"<Space>X",
		"<Space>X",
		true,
		keysequence.KeySequence{
			{termbox.KeySpace, ' ', 0},
			{0, 'X', 0},
		},
	},

	// Syntax errors
	{"", "", false, nil},
	{"<", "", false, nil},
	{"<space", "", false, nil},
	{"<>", "", false, nil},
	{"<C->", "", false, nil},
	{"<C-S->", "", false, nil},
	{"<X-m>", "", false, nil},
	{"<crap>", "", false, nil},
	{"<C-crap>", "", false, nil},
}

// Test that key sequences are correctly parsed.
func TestParser(t *testing.T) {
	for i, test := range parserTests {
		reader := strings.NewReader(test.input)
		scanner := lexer.NewScanner(reader)
		parser := keysequence.NewParser(scanner)

		t.Logf("Test %d: '%s'", i+1, test.input)

		seq, err := parser.ParseKeySequence()

		// Test success
		if test.success {
			assert.Nil(t, err, "Unexpected error when parsing '%s': %s", test.input, err)
		} else {
			assert.NotNil(t, err, "Expected error when parsing '%s'", test.input)
			continue
		}

		// Assert that names are converted back
		conv := keysequence.Format(seq)
		assert.Equal(t, test.output, conv, "Assert that reverse generated key sequence names are correct")

		// Assert that key definitions are equal
		assert.Equal(t, len(test.keyseq), len(seq), "Assert that key sequences have equal length")
		for k := range seq {
			t.Logf("Keyseq data in position %d: key=%d, rune='%s', mods=%d", k+1, seq[k].Key, string(seq[k].Ch), seq[k].Mod)
			if k >= len(test.keyseq) {
				continue
			}
			assert.Equal(t, test.keyseq[k].Key, seq[k].Key, "Assert that key event has equal Key in position %d", k+1)
			assert.Equal(t, test.keyseq[k].Ch, seq[k].Ch, "Assert that key event has equal Ch in position %d", k+1)
			assert.Equal(t, test.keyseq[k].Mod, seq[k].Mod, "Assert that key event has equal Mod in position %d", k+1)
		}
	}
}
