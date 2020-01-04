package keys_test

import (
	"github.com/ambientsound/pms/commands"
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
		err := sequencer.AddBind(keys.Binding{
			Sequence: seq,
			Command:  "foo bar",
			Context:  commands.GlobalContext,
		})
		assert.NoError(t, err)
	})

	t.Run("duplicate key bindings yield error", func(t *testing.T) {
		sequencer := keys.NewSequencer()
		err := sequencer.AddBind(keys.Binding{
			Sequence: seq,
			Command:  "foo bar",
			Context:  commands.GlobalContext,
		})
		assert.NoError(t, err)
		err = sequencer.AddBind(keys.Binding{
			Sequence: seq,
			Command:  "baz",
			Context:  commands.GlobalContext,
		})
		assert.Error(t, err)
	})

	t.Run("duplicate key bindings work if context is different", func(t *testing.T) {
		sequencer := keys.NewSequencer()
		err := sequencer.AddBind(keys.Binding{
			Sequence: seq,
			Command:  "foo bar",
			Context:  commands.GlobalContext,
		})
		assert.NoError(t, err)

		err = sequencer.AddBind(keys.Binding{
			Sequence: seq,
			Command:  "baz",
			Context:  commands.TracklistContext,
		})
		assert.NoError(t, err)

		// Test context preference way one
		contexts := []string{commands.TracklistContext, commands.GlobalContext}
		for i := range seq {
			sequencer.KeyInput(seq[i], contexts)
		}
		match := sequencer.Match(contexts)
		assert.NotNil(t, match)
		assert.Equal(t, "baz", match.Command)
		assert.Equal(t, commands.TracklistContext, match.Context)

		// Test context preference way two
		contexts = []string{commands.GlobalContext, commands.TracklistContext}
		for i := range seq {
			sequencer.KeyInput(seq[i], contexts)
		}
		match = sequencer.Match(contexts)
		assert.NotNil(t, match)
		assert.Equal(t, "foo bar", match.Command)
		assert.Equal(t, commands.GlobalContext, match.Context)
	})

	t.Run("matching a key binding", func(t *testing.T) {
		sequencer := keys.NewSequencer()
		contexts := []string{commands.GlobalContext}
		err := sequencer.AddBind(keys.Binding{
			Sequence: seq,
			Command:  "foo bar",
			Context:  commands.GlobalContext,
		})
		assert.NoError(t, err)

		assert.True(t, sequencer.KeyInput(seq[0], contexts))
		assert.Nil(t, sequencer.Match(contexts))
		assert.Equal(t, "a", sequencer.String())

		assert.True(t, sequencer.KeyInput(seq[1], contexts))
		assert.Nil(t, sequencer.Match(contexts))
		assert.Equal(t, "ab", sequencer.String())

		assert.True(t, sequencer.KeyInput(seq[2], contexts))
		assert.Equal(t, "abc", sequencer.String())
		match := sequencer.Match(contexts)

		assert.NotNil(t, match)
		assert.Equal(t, "foo bar", match.Command)
	})
}
