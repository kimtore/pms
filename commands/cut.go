package commands

import (
	"fmt"
	"github.com/ambientsound/pms/log"

	"github.com/ambientsound/pms/api"
)

// Cut removes songs from songlists.
type Cut struct {
	command
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
	list := cmd.api.Tracklist()
	// selection := list.Selection()
	indices := list.SelectionIndices()
	ln := len(indices)

	if ln == 0 {
		return fmt.Errorf("no tracks selected")
	}

	// Remove songs from list
	index := indices[0]
	err := list.RemoveIndices(indices)
	cmd.api.ListChanged()

	if err != nil {
		return err
	}

	log.Infof("%d fewer songs", ln)

	list.ClearSelection()
	list.SetCursor(index)

	// Place songs in clipboard
	// FIXME
	// clipboard := cmd.api.Db().Clipboard("default")
	// selection.Duplicate(clipboard)
	// console.Log("Cut %d tracks into clipboard", clipboard.Len())

	return nil
}
