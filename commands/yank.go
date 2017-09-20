package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
)

// Yank copies tracks from the songlist into the clipboard.
type Yank struct {
	newcommand
	api api.API
}

// NewYank returns Yank.
func NewYank(api api.API) Command {
	return &Yank{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Yank) Parse() error {
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Yank) Exec() error {
	list := cmd.api.Songlist()
	selection := list.Selection()
	indices := list.SelectionIndices()
	len := len(indices)

	if len == 0 {
		return fmt.Errorf("No tracks selected.")
	}

	// Place songs in clipboard
	clipboard := cmd.api.Db().Clipboard("default")
	selection.Duplicate(clipboard)

	// Print a message
	if len == 1 {
		cmd.api.Message("Yanked '%s'", selection.Song(0).StringTags["file"])
	} else {
		cmd.api.Message("%d tracks yanked to clipboard.", len)
	}

	// Clear selection and move cursor
	list.ClearSelection()
	list.MoveCursor(1)

	return nil
}
