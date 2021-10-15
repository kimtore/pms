package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
)

// Clears the queue queue
type Clear struct {
	newcommand
	api api.API
}

// NewClear returns Clear.
func NewClear(api api.API) Command {
	return &Clear{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Clear) Parse() error {
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Clear) Exec() error {
	if client := cmd.api.MpdClient(); client != nil {
		err := client.Clear()
		if err != nil {
			cmd.api.Message("clearing queue")
		}
		return err
	}
	return fmt.Errorf("Unable to clear: cannot communicate with MPD")
}
