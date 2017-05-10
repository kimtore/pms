package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Stop stops song playback in MPD.
type Stop struct {
	api api.API
}

func NewStop(api api.API) Command {
	return &Stop{
		api: api,
	}
}

func (cmd *Stop) Execute(class int, s string) error {
	switch class {
	case lexer.TokenEnd:
		if client := cmd.api.MpdClient(); client != nil {
			return client.Stop()
		}
		return fmt.Errorf("Unable to stop: cannot communicate with MPD")
	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}
}
