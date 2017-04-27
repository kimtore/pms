package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"

	"github.com/ambientsound/gompd/mpd"
)

// Stop plays songs in the MPD playlist.
type Stop struct {
	mpdClient func() *mpd.Client
}

func NewStop(mpdClient func() *mpd.Client) *Stop {
	return &Stop{mpdClient: mpdClient}
}

func (cmd *Stop) Reset() {
}

func (cmd *Stop) Execute(t lexer.Token) error {
	switch t.Class {
	case lexer.TokenEnd:
		if client := cmd.mpdClient(); client != nil {
			return client.Stop()
		}
		return fmt.Errorf("Unable to stop: cannot communicate with MPD")
	default:
		return fmt.Errorf("Unknown input '%s', expected END", t.String())
	}

	return nil
}
