package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"

	pms_mpd "github.com/ambientsound/pms/mpd"
)

// Pause toggles MPD play/paused state. If the player is stopped, Pause will
// attempt to start playback through the 'play' command instead.
type Pause struct {
	api api.API
}

func NewPause(api api.API) Command {
	return &Pause{
		api: api,
	}
}

func (cmd *Pause) Execute(class int, s string) error {
	switch class {
	case lexer.TokenEnd:
		client := cmd.api.MpdClient()
		if client == nil {
			return fmt.Errorf("Unable to toggle pause: cannot communicate with MPD")
		}
		status := cmd.api.PlayerStatus()
		switch status.State {
		case pms_mpd.StatePause:
			return client.Pause(false)
		case pms_mpd.StatePlay:
			return client.Pause(true)
		default:
			return client.Play(-1)
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}
}
