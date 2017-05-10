package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/mpd"
)

var stateStrings = map[string]string{
	mpd.StatePlay:    "|>",
	mpd.StatePause:   "||",
	mpd.StateStop:    "[]",
	mpd.StateUnknown: "??",
}

// State draws the current player state as an ASCII symbol.
type State struct {
	api api.API
}

// NewState returns State.
func NewState(a api.API, param string) Fragment {
	return &State{a}
}

// Text implements Fragment.
func (w *State) Text() (string, string) {
	playerStatus := w.api.PlayerStatus()
	return stateStrings[playerStatus.State], `state`
}
