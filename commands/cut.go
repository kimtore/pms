package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
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

	index := indices[0]
	err := list.RemoveIndices(indices)
	if err == nil {
		if len == 1 {
			cmd.api.Message("Cut out '%s'", selection.Song(0).StringTags["file"])
		} else {
			cmd.api.Message("%d fewer songs", len)
		}
		list.ClearSelection()
		list.SetCursor(index)
	}

	cmd.api.ListChanged()

	return nil
}
