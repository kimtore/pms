package keysequence

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
)

// KeySequence is an ordered sequence of keyboard events.
type KeySequence []*tcell.EventKey

// CompareKey compares two EventKey instances.
func CompareKey(a, b *tcell.EventKey) bool {
	if a.Modifiers() != b.Modifiers() || a.Key() != b.Key() {
		return false
	}
	// Runes don't have to match in case of a special key.
	if a.Key() != tcell.KeyRune {
		return true
	}
	return a.Rune() == b.Rune()
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

// FormatKey is similar to tcell.EventKey.Name(), which returns a printable
// value of a key stroke. Format formats it according to PMS' key binding syntax.
func FormatKey(ev *tcell.EventKey) string {
	s := ""
	m := []string{}
	mods := ev.Modifiers()

	// Add modifier keys
	if mods&tcell.ModShift != 0 {
		m = append(m, "Shift")
	}
	if mods&tcell.ModAlt != 0 {
		m = append(m, "Alt")
	}
	if mods&tcell.ModMeta != 0 {
		m = append(m, "Meta")
	}
	if mods&tcell.ModCtrl != 0 {
		m = append(m, "Ctrl")
	}

	// Check if the key already has a name. If not, use the correct rune. If
	// there is no matching rune, fall back to a question mark.
	ok := false
	key := ev.Key()
	if s, ok = tcell.KeyNames[key]; !ok {
		r := ev.Rune()
		if key == tcell.KeyRune {
			if r == ' ' {
				s = "<Space>"
			} else {
				s = string(ev.Rune())
			}
		} else {
			s = fmt.Sprintf("<%d,%d>", key, int(r))
		}
	}

	// Append any modifier prefixes.
	if len(m) != 0 {
		if mods&tcell.ModCtrl != 0 && strings.HasPrefix(s, "Ctrl-") {
			s = s[5:]
		}
		return fmt.Sprintf("<%s-%s>", strings.Join(m, "-"), s)
	}
	return s
}

// Format reverses a parsed key sequence into its string representation.
func Format(seq KeySequence) string {
	s := make([]string, len(seq))
	for i := range seq {
		s[i] = FormatKey(seq[i])
	}
	return strings.Join(s, "")
}
