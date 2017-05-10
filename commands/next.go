package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Next switches to the next song in MPD's queue.
type Next struct {
	api api.API
}

func NewNext(api api.API) Command {
	return &Next{
		api: api,
	}
}

func (cmd *Next) Execute(class int, s string) error {
	switch class {
	case lexer.TokenEnd:
		client := cmd.api.MpdClient()
		if client == nil {
			return fmt.Errorf("Unable to play next song: cannot communicate with MPD")
		}
		return client.Next()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}
}
