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
		selection := list.SelectionIndices()

		if len(selection) == 0 {
			return fmt.Errorf("No song selected, cannot remove without any parameters.")
		}

		index := selection[0]
		err = list.RemoveIndices(selection)
		if err == nil {
			list.ClearSelection()
			list.SetCursor(index)
		}

		cmd.api.ListChanged()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}
