package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/utils"
)

// Time draws the current volume.
type Time struct {
	api api.API
}

func NewTime(a api.API, param string) Fragment {
	return &Time{a}
}

func (w *Time) Text() (string, string) {
	playerStatus := w.api.PlayerStatus()
	return utils.TimeString(playerStatus.Time), `time`
}
