package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Previous switches to the previous song in MPD's queue.
type Previous struct {
	api api.API
}

func NewPrevious(api api.API) Command {
	return &Previous{
		api: api,
	}
}

func (cmd *Previous) Execute(class int, s string) error {
	switch class {
	case lexer.TokenEnd:
		client := cmd.api.MpdClient()
		if client == nil {
			return fmt.Errorf("Unable to play previous song: cannot communicate with MPD")
		}
		return client.Previous()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}
}
