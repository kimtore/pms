package term

import (
	termbox "github.com/nsf/termbox-go"
)

type Modifier uint8

const (
	ModCtrl Modifier = 1 << iota
	ModAlt
	ModShift
	ModMeta
)

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

	k = convertCtrlKey(k)

	return k
}

// convertCtrlKey sets the rune and modifier members if a keypress is Ctrl+A through Ctrl+Z.
func convertCtrlKey(k KeyPress) KeyPress {
	if k.Key >= termbox.KeyCtrlA && k.Key <= termbox.KeyCtrlZ {
		k.Ch = rune(k.Key + 96)
		k.Mod |= ModCtrl
	}
	return k
}
