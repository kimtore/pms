package keys_test

import (
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/keysequence"
	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequencer(t *testing.T) {

	seq := keysequence.KeySequence{
		tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
		tcell.NewEventKey(tcell.KeyRune, 'b', tcell.ModNone),
		tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
	}

	t.Run("adding key bindings yields no error", func(t *testing.T) {
		sequencer := keys.NewSequencer()
		err := sequencer.AddBind(seq, "foo bar")
		assert.NoError(t, err)
	})

	t.Run("duplicate key bindings yield error", func(t *testing.T) {
		sequencer := keys.NewSequencer()
		err := sequencer.AddBind(seq, "foo bar")
		assert.NoError(t, err)
		err = sequencer.AddBind(seq, "baz")
		assert.Error(t, err)
	})

	t.Run("matching a key binding", func(t *testing.T) {
		sequencer := keys.NewSequencer()
		err := sequencer.AddBind(seq, "foo bar")
		assert.NoError(t, err)

		assert.True(t, sequencer.KeyInput(seq[0]))
		assert.Nil(t, sequencer.Match())
		assert.Equal(t, "a", sequencer.String())

		assert.True(t, sequencer.KeyInput(seq[1]))
		assert.Nil(t, sequencer.Match())
		assert.Equal(t, "ab", sequencer.String())

		assert.True(t, sequencer.KeyInput(seq[2]))
		assert.Equal(t, "abc", sequencer.String())
		match := sequencer.Match()

		assert.NotNil(t, match)
		assert.Equal(t, "foo bar", match.Command)
	})
}
