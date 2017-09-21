package term

import (
	"fmt"
	"strings"
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

var KeyNames = map[string]termbox.Key{
	"Backspace":  termbox.KeyBackspace,
	"Backspace2": termbox.KeyBackspace2,
	"Delete":     termbox.KeyDelete,
	"Down":       termbox.KeyArrowDown,
	"End":        termbox.KeyEnd,
	"Enter":      termbox.KeyEnter,
	"Esc":        termbox.KeyEsc,
	"F1":         termbox.KeyF1,
	"F2":         termbox.KeyF2,
	"F3":         termbox.KeyF3,
	"F4":         termbox.KeyF4,
	"F5":         termbox.KeyF5,
	"F6":         termbox.KeyF6,
	"F7":         termbox.KeyF7,
	"F8":         termbox.KeyF8,
	"F9":         termbox.KeyF9,
	"F10":        termbox.KeyF10,
	"F11":        termbox.KeyF11,
	"F12":        termbox.KeyF12,
	"Home":       termbox.KeyHome,
	"Insert":     termbox.KeyInsert,
	"Left":       termbox.KeyArrowLeft,
	"PgDn":       termbox.KeyPgdn,
	"PgUp":       termbox.KeyPgup,
	"Right":      termbox.KeyArrowRight,
	"Space":      termbox.KeySpace,
	"Tab":        termbox.KeyTab,
	"Up":         termbox.KeyArrowUp,
}

func (k *KeyPress) keyName() (string, error) {
	for key, val := range KeyNames {
		if k.Key == val {
			return key, nil
		}
	}
	return "", fmt.Errorf("No name for this key")
}

// FormatKey is similar to tcell.EventKey.Name(), which returns a printable
// value of a key stroke. Format formats it according to PMS' key binding syntax.
func (k *KeyPress) Name() string {
	s := ""
	m := []string{}

	// Add modifier keys
	if k.Mod&ModShift != 0 {
		m = append(m, "Shift")
	}
	if k.Mod&ModAlt != 0 {
		m = append(m, "Alt")
	}
	if k.Mod&ModMeta != 0 {
		m = append(m, "Meta")
	}
	if k.Mod&ModCtrl != 0 {
		m = append(m, "Ctrl")
	}

	// Check if the key has a pre-defined name. If not, use the correct rune. If there is no matching rune, fall back to a question mark.
	s, err := k.keyName()
	if err != nil {
		if k.Ch == 0 {
			s = fmt.Sprintf("<%d,%d>", k.Key, int(k.Ch))
		} else {
			// Fall back to using the rune. If Ctrl is held, the rune should be uppercased.
			if k.Mod&ModCtrl != 0 {
				s = string(unicode.ToUpper(k.Ch))
			} else {
				s = string(k.Ch)
			}
		}
	}

	// Append any modifier prefixes.
	if len(m) != 0 {
		s = fmt.Sprintf("%s-%s", strings.Join(m, "-"), s)
	}

	if len(s) > 1 {
		s = fmt.Sprintf("<%s>", s)
	}

	return s
}
