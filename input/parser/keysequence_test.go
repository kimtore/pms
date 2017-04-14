package parser_test

import (
	"testing"

	"github.com/ambientsound/pms/input/parser"
	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapParser(t *testing.T) {
	input := "<C-d>tr<o>ll<F1>"
	tokens := parser.KeySequenceToken{}

	err := tokens.Parse([]rune(input))
	assert.Nil(t, err)

	require.Equal(t, 9, len(tokens.Sequence))

	assert.Equal(t, tcell.KeyCtrlD, tokens.Sequence[0].Key)

	for i := 1; i < len(tokens.Sequence)-1; i++ {
		assert.Equal(t, tcell.KeyRune, tokens.Sequence[i].Key)
	}
	assert.Equal(t, 't', tokens.Sequence[1].Rune)
	assert.Equal(t, 'r', tokens.Sequence[2].Rune)
	assert.Equal(t, '<', tokens.Sequence[3].Rune)
	assert.Equal(t, 'o', tokens.Sequence[4].Rune)
	assert.Equal(t, '>', tokens.Sequence[5].Rune)
	assert.Equal(t, 'l', tokens.Sequence[6].Rune)
	assert.Equal(t, 'l', tokens.Sequence[7].Rune)

	assert.Equal(t, tcell.KeyF1, tokens.Sequence[8].Key)
}
