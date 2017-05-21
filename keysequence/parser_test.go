package keysequence_test

import (
	"strings"
	"testing"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/keysequence"
	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
)

var parserTests = []struct {
	input   string
	success bool
	keyseq  keysequence.KeySequence
}{
	{
		"abc",
		true,
		keysequence.KeySequence{
			tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyRune, 'b', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
		},
	},
	{
		"<C-c>x<S-m-A-f1>f",
		true,
		keysequence.KeySequence{
			tcell.NewEventKey(tcell.KeyCtrlC, rune(tcell.KeyCtrlC), tcell.ModCtrl),
			tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyF1, 0, tcell.ModShift|tcell.ModMeta|tcell.ModAlt),
			tcell.NewEventKey(tcell.KeyRune, 'f', tcell.ModNone),
		},
	},
	{
		"<a-x>",
		true,
		keysequence.KeySequence{
			tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModAlt),
		},
	},
	{
		"<c-x>x",
		true,
		keysequence.KeySequence{
			tcell.NewEventKey(tcell.KeyCtrlX, rune(tcell.KeyCtrlX), tcell.ModCtrl),
			tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone),
		},
	},
	{
		"<Space>",
		true,
		keysequence.KeySequence{
			tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone),
		},
	},
	{
		"<Space>X",
		true,
		keysequence.KeySequence{
			tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone),
			tcell.NewEventKey(tcell.KeyRune, 'X', tcell.ModNone),
		},
	},

	// Syntax errors
	{"", false, nil},
	{"<", false, nil},
	{"<space", false, nil},
	{"<>", false, nil},
	{"<C->", false, nil},
	{"<C-S->", false, nil},
	{"<X-m>", false, nil},
	{"<crap>", false, nil},
	{"<C-crap>", false, nil},
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

		// Assert that key definitions are equal
		assert.Equal(t, len(test.keyseq), len(seq), "Assert that key sequences have equal length")
		for k := range seq {
			t.Logf("Keyseq data in position %d: key=%d, rune='%s', mods=%d", k+1, seq[k].Key(), string(seq[k].Rune()), seq[k].Modifiers())
			if k >= len(test.keyseq) {
				continue
			}
			assert.Equal(t, test.keyseq[k].Key(), seq[k].Key(), "Assert that key event has equal Key() in position %d", k+1)
			assert.Equal(t, test.keyseq[k].Rune(), seq[k].Rune(), "Assert that key event has equal Rune() in position %d", k+1)
			assert.Equal(t, test.keyseq[k].Modifiers(), seq[k].Modifiers(), "Assert that key event has equal Modifiers() in position %d", k+1)
		}
	}
}
