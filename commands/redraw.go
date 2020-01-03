package commands

import (
	"github.com/ambientsound/pms/api"
)

// Redraw forcefully redraws the screen.
type Redraw struct {
	command
	api api.API
}

func NewRedraw(api api.API) Command {
	return &Redraw{
		api: api,
	}
}

func (cmd *Redraw) Parse() error {
	return cmd.ParseEnd()
}

func (cmd *Redraw) Exec() error {
	cmd.api.UI().Refresh()
	return nil
}
