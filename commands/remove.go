package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Remove removes songs from songlists.
type Remove struct {
	api api.API
}

func NewRemove(api api.API) Command {
	return &Remove{
		api: api,
	}
}

func (cmd *Remove) Execute(class int, s string) error {
	var err error

	switch class {
	case lexer.TokenEnd:
		list := cmd.api.Songlist()
		selection := list.Selection()
		indices := list.SelectionIndices()
		len := len(indices)

		if len == 0 {
			return fmt.Errorf("No song selected, cannot remove without any parameters.")
		}

		index := indices[0]
		err = list.RemoveIndices(indices)
		if err == nil {
			if len == 1 {
				cmd.api.Message("Removed '%s'", selection.Song(0).StringTags["file"])
			} else {
				cmd.api.Message("%d fewer songs", len)
			}
			list.ClearSelection()
			list.SetCursor(index)
		}

		cmd.api.ListChanged()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}
