package topbar

import (
"bytes"

"github.com/ambientsound/pms/api"
)

// Mode draws the four player modes as single characters.
type Mode struct {
	api api.API
}

// NewMode returns Mode.
func NewMode(a api.API, param string) Fragment {
	return &Mode{a}
}

// Text implements Fragment.
func (w *Mode) Text() (string, string) {
	var buf bytes.Buffer
	playerStatus := w.api.PlayerStatus()

	buf.WriteRune(w.statusRune('z', playerStatus.ShuffleState))
	buf.WriteRune(w.statusRune('s', playerStatus.RepeatState == "track"))
	buf.WriteRune(w.statusRune('r', playerStatus.RepeatState != "off"))

	return buf.String(), `switches`
}

func (w *Mode) statusRune(r rune, val bool) rune {
	if val {
		return r
	}
	return '-'
}
