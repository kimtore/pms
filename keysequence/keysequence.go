package keysequence

import (
	"strings"

	"github.com/ambientsound/pms/term"
)

// KeySequence is an ordered sequence of keyboard events.
type KeySequence []term.KeyPress

// CompareKey compares two EventKey instances.
func CompareKey(a, b term.KeyPress) bool {
	if a.Mod != b.Mod || a.Key != b.Key {
		return false
	}
	// Runes don't have to match in case of a special key.
	if a.Key != 0 {
		return true
	}
	return a.Ch == b.Ch
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

// Format reverses a parsed key sequence into its string representation.
func Format(seq KeySequence) string {
	s := make([]string, len(seq))
	for i := range seq {
		s[i] = seq[i].Name()
	}
	return strings.Join(s, "")
}
