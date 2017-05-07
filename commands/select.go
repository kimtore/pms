package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
)

// Select manipulates song selection within a songlist.
type Select struct {
	api      API
	toggle   bool
	finished bool
}

func NewSelect(api API) Command {
	return &Select{
		api: api,
	}
}

func (cmd *Select) Execute(t lexer.Token) error {
	var err error

	s := t.String()
	songlistWidget := cmd.api.SonglistWidget()
	list := songlistWidget.List()

	switch t.Class {

	case lexer.TokenIdentifier:
		if cmd.finished {
			return fmt.Errorf("Unexpected '%s', expected END", s)
		}
		switch s {
		case "toggle":
			cmd.toggle = true
		default:
			return fmt.Errorf("Unexpected '%s', expected identifier", s)
		}
		cmd.finished = true

	case lexer.TokenEnd:
		if !cmd.finished {
			return fmt.Errorf("Unexpected END, expected identifier")
		}

		switch {
		case list.HasVisualSelection():
			list.CommitVisualSelection()
			songlistWidget.DisableVisualSelection()

		default:
			index := songlistWidget.Cursor()
			selected := songlistWidget.List().Selected(index)
			songlistWidget.List().SetSelected(index, !selected)
		}

		songlistWidget.MoveCursor(1)

	default:
		return fmt.Errorf("Unexpected '%s', expected END", s)
	}

	return err
}
