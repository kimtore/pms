package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/player"
)

var stateStrings = map[string]string{
	player.StatePlay:    "|>",
	player.StatePause:   "||",
	player.StateStop:    "[]",
	player.StateUnknown: "??",
}

var stateUnicodes = map[string]string{
	player.StatePlay:    "\u25b6",
	player.StatePause:   "\u23f8",
	player.StateStop:    "\u23f9",
	player.StateUnknown: "\u2bd1",
}

// State draws the current player state as an ASCII symbol.
type State struct {
	api   api.API
	table map[string]string
}

// NewState returns State.
func NewState(a api.API, param string) Fragment {
	table := stateStrings
	if param == "unicode" {
		table = stateUnicodes
	}
	return &State{a, table}
}

// Text implements Fragment.
func (w *State) Text() (string, string) {
	playerStatus := w.api.PlayerStatus()
	return w.table[playerStatus.State()], `state`
}
