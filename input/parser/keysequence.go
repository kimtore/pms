package parser

import (
	"strings"

	"github.com/gdamore/tcell"
)

type KeyEvent struct {
	Key  tcell.Key
	Rune rune
}

func (a KeyEvent) Equals(b KeyEvent) bool {
	if a.Key != b.Key {
		return false
	}
	if a.Key == tcell.KeyRune {
		if a.Rune != b.Rune {
			return false
		}
	}
	return true
}

type KeyEvents []KeyEvent

func (a KeyEvents) Equals(b KeyEvents) bool {
	if len(a) != len(b) {
		return false
	}
	return a.StartsWith(b)
}

func (a KeyEvents) StartsWith(b KeyEvents) bool {
	if len(a) < len(b) {
		return false
	}
	for i := range b {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}

type KeySequenceToken struct {
	Sequence KeyEvents
}

func (t *KeySequenceToken) addRunes(runes []rune) {
	for _, r := range runes {
		t.Sequence = append(t.Sequence, KeyEvent{Key: tcell.KeyRune, Rune: r})
	}
}

// Parse parses a sequence of keystrokes defined as a string, and creates a
// slice of KeyEvent structs, representing individual keystrokes.
func (t *KeySequenceToken) Parse(runes []rune) error {
	t.Sequence = make(KeyEvents, 0)

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
			if _, ok := keyNames[s]; ok {
				t.Sequence = append(t.Sequence, keyNames[s])
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
