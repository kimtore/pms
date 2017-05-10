package topbar

import (
	"bytes"

	"github.com/ambientsound/pms/api"
)

// Mode draws the current volume.
type Mode struct {
	api api.API
}

func NewMode(a api.API, param string) Fragment {
	return &Mode{a}
}

func (w *Mode) Text() (string, string) {
	var buf bytes.Buffer
	playerStatus := w.api.PlayerStatus()

	buf.WriteRune(w.statusRune('c', playerStatus.Consume))
	buf.WriteRune(w.statusRune('z', playerStatus.Random))
	buf.WriteRune(w.statusRune('s', playerStatus.Single))
	buf.WriteRune(w.statusRune('r', playerStatus.Repeat))

	return buf.String(), `switches`
}

func (w *Mode) statusRune(r rune, val bool) rune {
	if val {
		return r
	}
	return '-'
}
