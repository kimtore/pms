package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
)

// Remove removes songs from songlists.
type Remove struct {
	api API
}

func NewRemove(api API) Command {
	return &Remove{
		api: api,
	}
}

func (cmd *Remove) Execute(t lexer.Token) error {
	var err error

	switch t.Class {
	case lexer.TokenEnd:
		songlistWidget := cmd.api.SonglistWidget()
		list := songlistWidget.Songlist()
		selection := songlistWidget.List().SelectionIndices()

		if len(selection) == 0 {
			return fmt.Errorf("No song selected, cannot remove without any parameters.")
		}

		index := selection[0]
		err = list.RemoveIndices(selection)
		if err == nil {
			songlistWidget.ClearSelection()
			songlistWidget.SetCursor(index)
		}

		cmd.api.ListChanged()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return err
}
