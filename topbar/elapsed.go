package topbar

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/utils"
)

// Elapsed draws the current song's elapsed time.
type Elapsed struct {
	api api.API
	f   func() (string, string)
}

// NewElapsed returns Elapsed.
func NewElapsed(a api.API, param string) Fragment {
	elapsed := &Elapsed{a, nil}
	switch param {
	case `percentage`:
		elapsed.f = elapsed.textPercentage
	default:
		elapsed.f = elapsed.textTime
	}
	return elapsed
}

// Text implements Fragment.
func (w *Elapsed) Text() (string, string) {
	return w.f()
}

func (w *Elapsed) textTime() (string, string) {
	playerStatus := w.api.PlayerStatus()
	return utils.TimeString(int(playerStatus.Elapsed)), `elapsedTime`
}

func (w *Elapsed) textPercentage() (string, string) {
	playerStatus := w.api.PlayerStatus()
	return fmt.Sprintf("%d", int(playerStatus.ElapsedPercentage)), `elapsedPercentage`
}
