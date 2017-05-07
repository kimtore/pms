package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
)

// Stop stops song playback in MPD.
type Stop struct {
	api API
}

func NewStop(api API) Command {
	return &Stop{
		api: api,
	}
}

func (cmd *Stop) Execute(t lexer.Token) error {
	switch t.Class {
	case lexer.TokenEnd:
		if client := cmd.api.MpdClient(); client != nil {
			return client.Stop()
		}
		return fmt.Errorf("Unable to stop: cannot communicate with MPD")
	default:
		return fmt.Errorf("Unknown input '%s', expected END", t.String())
	}

	return nil
}
