package commands

import (
	"github.com/ambientsound/pms/api"
	termbox "github.com/nsf/termbox-go"
)

// Redraw tries to synchronize the terminal backbuffer and the screen.
type Redraw struct {
	newcommand
	api api.API
}

// NewRedraw returns Redraw.
func NewRedraw(api api.API) Command {
	return &Redraw{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Redraw) Parse() error {
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Redraw) Exec() error {
	// FIXME: set all models as dirty
	return termbox.Sync()
}
