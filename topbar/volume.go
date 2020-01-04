package topbar

import (
	"fmt"

	"github.com/ambientsound/pms/api"
)

// Volume draws the current volume.
type Volume struct {
	api api.API
}

// NewVolume returns Volume.
func NewVolume(a api.API, param string) Fragment {
	return &Volume{a}
}

// Text implements Fragment.
func (w *Volume) Text() (string, string) {
	playerStatus := w.api.PlayerStatus()
	switch {
	case playerStatus.Device.Volume < 0:
		return `!VOL!`, `mute`
	case playerStatus.Device.Volume == 0:
		return `MUTE`, `mute`
	default:
		text := fmt.Sprintf("%d%%", playerStatus.Device.Volume)
		return text, `volume`
	}
}
