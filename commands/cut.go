package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
)

// Cut removes songs from songlists.
type Cut struct {
	newcommand
	api api.API
}

// NewCut returns Cut.
func NewCut(api api.API) Command {
	return &Cut{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Cut) Parse() error {
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Cut) Exec() error {
	list := cmd.api.Songlist()
	selection := list.Selection()
	indices := list.SelectionIndices()
	len := len(indices)

	if len == 0 {
		return fmt.Errorf("No tracks selected, cannot remove without any parameters.")
	}

	// Remove songs from list
	index := indices[0]
	err := list.RemoveIndices(indices)
	cmd.api.ListChanged()

	if err != nil {
		return err
	}

	if len == 1 {
		cmd.api.Message("Cut out '%s'", selection.Song(0).StringTags["file"])
	} else {
		cmd.api.Message("%d fewer songs", len)
	}
	list.ClearSelection()
	list.SetCursor(index)

	// Place songs in clipboard
	clipboard := cmd.api.Db().Clipboard("default")
	selection.Duplicate(clipboard)
	console.Log("Cut %d tracks into clipboard", clipboard.Len())

	return nil
}
