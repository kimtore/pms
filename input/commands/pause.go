package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"

	"github.com/ambientsound/gompd/mpd"
	pms_mpd "github.com/ambientsound/pms/mpd"
)

// Pause toggles MPD play/paused state. If the player is stopped, Pause will
// attempt to start playback through the 'play' command instead.
type Pause struct {
	mpdClient func() *mpd.Client
	mpdStatus func() pms_mpd.PlayerStatus
}

func NewPause(mpdClient func() *mpd.Client, mpdStatus func() pms_mpd.PlayerStatus) *Pause {
	return &Pause{mpdClient: mpdClient, mpdStatus: mpdStatus}
}

func (cmd *Pause) Reset() {
}

func (cmd *Pause) Execute(t lexer.Token) error {
	switch t.Class {
	case lexer.TokenEnd:
		client := cmd.mpdClient()
		if client == nil {
			return fmt.Errorf("Unable to toggle pause: cannot communicate with MPD")
		}
		status := cmd.mpdStatus()
		switch status.State {
		case pms_mpd.StatePause:
			return client.Pause(false)
		case pms_mpd.StatePlay:
			return client.Pause(true)
		default:
			return client.Play(-1)
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", t.String())
	}

	return nil
}
