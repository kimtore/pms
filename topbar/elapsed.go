package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/utils"
)

// Elapsed draws the current song's elapsed time.
type Elapsed struct {
	api api.API
}

// NewElapsed returns Elapsed.
func NewElapsed(a api.API, param string) Fragment {
	return &Elapsed{a}
}

// Text implements Fragment.
func (w *Elapsed) Text() (string, string) {
	playerStatus := w.api.PlayerStatus()
	return utils.TimeString(int(playerStatus.Elapsed)), `elapsed`
}
