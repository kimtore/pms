package keys_test

import (
	"testing"

	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/input/parser"
	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSequencer(t *testing.T) {
	s := keys.NewSequencer()

	tok := parser.KeySequenceToken{}
	tok.Parse([]rune("abc"))
	s.AddBind(tok.Sequence, "foobar")

	tok = parser.KeySequenceToken{}
	tok.Parse([]rune("abcde"))
	require.NotNil(t, s.AddBind(tok.Sequence, "blah"))

	tok = parser.KeySequenceToken{}
	tok.Parse([]rune("<C-a><space>"))
	s.AddBind(tok.Sequence, "baz")

	assert.True(t, s.KeyInput(parser.KeyEvent{Key: tcell.KeyRune, Rune: 'a'}))
	assert.True(t, s.KeyInput(parser.KeyEvent{Key: tcell.KeyRune, Rune: 'b'}))
	assert.True(t, s.KeyInput(parser.KeyEvent{Key: tcell.KeyRune, Rune: 'c'}))

	in := s.Match()
	require.NotNil(t, in)
	assert.Equal(t, "foobar", in.Command)

	assert.False(t, s.KeyInput(parser.KeyEvent{Key: tcell.KeyRune, Rune: 'x'}))
	assert.False(t, s.KeyInput(parser.KeyEvent{Key: tcell.KeyRune, Rune: 'y'}))

	assert.True(t, s.KeyInput(parser.KeyEvent{Key: tcell.KeyCtrlA}))
	assert.True(t, s.KeyInput(parser.KeyEvent{Key: tcell.KeyRune, Rune: ' '}))
	in = s.Match()
	require.NotNil(t, in)
	assert.Equal(t, "baz", in.Command)
}
