package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"

	"github.com/ambientsound/gompd/mpd"
)

// Previous switches to the previous song in MPD's queue.
type Previous struct {
	mpdClient func() *mpd.Client
}

func NewPrevious(mpdClient func() *mpd.Client) *Previous {
	return &Previous{mpdClient: mpdClient}
}

func (cmd *Previous) Reset() {
}

func (cmd *Previous) Execute(t lexer.Token) error {
	switch t.Class {
	case lexer.TokenEnd:
		client := cmd.mpdClient()
		if client == nil {
			return fmt.Errorf("Unable to play previous song: cannot communicate with MPD")
		}
		return client.Previous()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", t.String())
	}

	return nil
}
