package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
)

// Update updates MPD database.
type Update struct {
	newcommand
	api api.API
}

// NewUpdate returns Update.
func NewUpdate(api api.API) Command {
	return &Update{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Update) Parse() error {
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Update) Exec() error {
	if client := cmd.api.MpdClient(); client != nil {
		jobID, err := client.Update("")
		if err != nil {
			cmd.api.Message("Updating, jobID is %d", jobID)
		}
		return err
	}
	return fmt.Errorf("Unable to update database: cannot communicate with MPD")
}
