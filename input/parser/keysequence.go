package parser

import (
	"strings"

	"github.com/gdamore/tcell"
)

type KeyEvent struct {
	Key  tcell.Key
	Rune rune
}

func (k KeyEvent) Equals(x KeyEvent) bool {
	if k.Key != x.Key {
		return false
	}
	if k.Key == tcell.KeyRune {
		if k.Rune != x.Rune {
			return false
		}
	}
	return true
}

func (k KeyEvent) String() string {
	if k.Key == tcell.KeyRune {
		return string(k.Rune)
	}
	for i := range keyNames {
		if k.Equals(keyNames[i]) {
			return "<" + i + ">"
		}
	}
	return "<UNKNOWN>"
}

type KeyEvents []KeyEvent

func (k KeyEvents) Equals(x KeyEvents) bool {
	if len(k) != len(x) {
		return false
	}
	return k.StartsWith(x)
}

func (k KeyEvents) StartsWith(x KeyEvents) bool {
	if len(k) < len(x) {
		return false
	}
	for i := range x {
		if !k[i].Equals(x[i]) {
			return false
		}
	}
	return true
}

func (k KeyEvents) String() string {
	s := make([]string, 0)
	for i := range k {
		s = append(s, k[i].String())
	}
	return strings.Join(s, "")
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

	parseSpecial := false
	specialCharacters := make([]rune, 0)

	for _, r := range runes {
		switch r {
		case '<':
			if parseSpecial {
				// If already parsing specials, assume that every key up to
				// this point is literal, and add them to the key sequence
				t.addRunes(specialCharacters)
				specialCharacters = make([]rune, 0)
			}
			specialCharacters = append(specialCharacters, r)
			parseSpecial = true
			continue
		case '>':
			if !parseSpecial {
				break
			}
			s := strings.ToLower(string(specialCharacters[1:]))
			if _, ok := keyNames[s]; ok {
				t.Sequence = append(t.Sequence, keyNames[s])
			} else {
				specialCharacters = append(specialCharacters, r)
				t.addRunes(specialCharacters)
			}
			specialCharacters = make([]rune, 0)
			parseSpecial = false
			continue
		}
		if parseSpecial {
			specialCharacters = append(specialCharacters, r)
		} else {
			t.Sequence = append(t.Sequence, KeyEvent{Key: tcell.KeyRune, Rune: r})
		}
	}

	return nil
}
