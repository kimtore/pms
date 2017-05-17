package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
)

// Stop stops song playback in MPD.
type Stop struct {
	newcommand
	api api.API
}

// NewStop returns Stop.
func NewStop(api api.API) Command {
	return &Stop{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Stop) Parse() error {
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Stop) Exec() error {
	if client := cmd.api.MpdClient(); client != nil {
		return client.Stop()
	}
	return fmt.Errorf("Unable to stop: cannot communicate with MPD")
}
