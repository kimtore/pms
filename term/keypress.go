package term

import (
	"fmt"
	"strings"
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

// Modifier represents a key modifier such as ctrl, alt, shift, and meta.
type Modifier uint8

// Keyboard modifiers. Currently, only Ctrl and Alt are detected.
const (
	ModCtrl Modifier = 1 << iota
	ModAlt
	ModShift
	ModMeta
)

// KeyNames holds a mapping from string names to termbox key constants.
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

// lowerKeyNames is like KeyNames, but all keys are in lowercase.
var lowerKeyNames = map[string]termbox.Key{}

// KeyPress represents a single keypress.
type KeyPress struct {
	Key termbox.Key
	Ch  rune
	Mod Modifier
}

// ParseKey converts a termbox keypress into a PMS keypress.
func ParseKey(te termbox.Event) KeyPress {
	k := KeyPress{
		Key: te.Key,
		Ch:  te.Ch,
	}

	if te.Mod == termbox.ModAlt {
		k.Mod = ModAlt
	}

	k = k.ConvertCtrlKey()

	if k.Key == termbox.KeySpace {
		k.Ch = ' '
	}

	return k
}

// Key takes a case insensitive string value and returns a termbox key constant.
func Key(key string) (termbox.Key, error) {
	if len(lowerKeyNames) == 0 {
		for k, v := range KeyNames {
			lowerKeyNames[strings.ToLower(k)] = v
		}
	}
	v, ok := lowerKeyNames[strings.ToLower(key)]
	if !ok {
		return 0, fmt.Errorf("No such key: %s", key)
	}
	return v, nil
}

// ConvertCtrlKey returns a new KeyPress with the rune and modifier members set
// accordingly if a keypress is Ctrl+A through Ctrl+Z.
func (k KeyPress) ConvertCtrlKey() KeyPress {
	if k.Key >= termbox.KeyCtrlA && k.Key <= termbox.KeyCtrlZ {
		k.Ch = rune(k.Key + 96)
		k.Mod |= ModCtrl
	}
	return k
}

// constName returns a string representation of a termbox key constant.
func (k KeyPress) constName() (string, error) {
	for key, val := range KeyNames {
		if k.Key == val {
			return key, nil
		}
	}
	return "", fmt.Errorf("No name for this key")
}

// Name returns a human-readable, canonical representation of a key stroke.
func (k KeyPress) Name() string {
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
	s, err := k.constName()
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
