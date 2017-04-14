package parser

import (
	"strings"

	"github.com/gdamore/tcell"
)

type KeyEvent struct {
	Key  tcell.Key
	Rune rune
}

type KeySequenceToken struct {
	Sequence []KeyEvent
}

func (t *KeySequenceToken) addRunes(runes []rune) {
	for _, r := range runes {
		t.Sequence = append(t.Sequence, KeyEvent{Key: tcell.KeyRune, Rune: r})
	}
}

// Parse parses a sequence of keystrokes defined as a string, and creates a
// slice of KeyEvent structs, representing individual keystrokes.
func (t *KeySequenceToken) Parse(runes []rune) error {
	t.Sequence = make([]KeyEvent, 0)

	parse_special := false
	special_characters := make([]rune, 0)

	for _, r := range runes {
		switch r {
		case '<':
			if parse_special {
				// If already parsing specials, assume that every key up to
				// this point is literal, and add them to the key sequence
				t.addRunes(special_characters)
				special_characters = make([]rune, 0)
			}
			special_characters = append(special_characters, r)
			parse_special = true
			continue
		case '>':
			if !parse_special {
				break
			}
			s := strings.ToLower(string(special_characters[1:]))
			if key, ok := keyNames[s]; ok {
				t.Sequence = append(t.Sequence, KeyEvent{Key: key})
			} else {
				special_characters = append(special_characters, r)
				t.addRunes(special_characters)
			}
			special_characters = make([]rune, 0)
			parse_special = false
			continue
		}
		if parse_special {
			special_characters = append(special_characters, r)
		} else {
			t.Sequence = append(t.Sequence, KeyEvent{Key: tcell.KeyRune, Rune: r})
		}
	}

	return nil
}
