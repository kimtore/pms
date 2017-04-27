package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"

	"github.com/ambientsound/gompd/mpd"
)

// Next toggles MPD play/paused state. If the player is stopped, Next will
// attempt to start playback through the 'play' command instead.
type Next struct {
	mpdClient func() *mpd.Client
}

func NewNext(mpdClient func() *mpd.Client) *Next {
	return &Next{mpdClient: mpdClient}
}

func (cmd *Next) Reset() {
}

func (cmd *Next) Execute(t lexer.Token) error {
	switch t.Class {
	case lexer.TokenEnd:
		client := cmd.mpdClient()
		if client == nil {
			return fmt.Errorf("Unable to play next song: cannot communicate with MPD")
		}
		return client.Next()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", t.String())
	}

	return nil
}
