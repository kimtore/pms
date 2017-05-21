package keysequence

import (
	"github.com/gdamore/tcell"
)

// KeySequence is an ordered sequence of keyboard events.
type KeySequence []*tcell.EventKey

// CompareKey compares two EventKey instances.
func CompareKey(a, b *tcell.EventKey) bool {
	if a.Modifiers() != b.Modifiers() {
		return false
	}
	if a.Rune() != b.Rune() {
		return false
	}
	return a.Key() == b.Key()
}

// Compare compares two KeySequence instances.
func Compare(a, b KeySequence) bool {
	if len(a) != len(b) {
		return false
	}
	return StartsWith(a, b)
}

// StartsWith return true if a starts with b.
func StartsWith(a, b KeySequence) bool {
	if len(b) > len(a) {
		return false
	}
	for i := range b {
		if !CompareKey(a[i], b[i]) {
			return false
		}
	}
	return true
}
