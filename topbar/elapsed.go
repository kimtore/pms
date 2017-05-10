package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/utils"
)

// Elapsed draws the current volume.
type Elapsed struct {
	api api.API
}

func NewElapsed(a api.API, param string) Fragment {
	return &Elapsed{a}
}

func (w *Elapsed) Text() (string, string) {
	playerStatus := w.api.PlayerStatus()
	return utils.TimeString(int(playerStatus.Elapsed)), `elapsed`
}
