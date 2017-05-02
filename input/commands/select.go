package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/widgets"
)

// Select manipulates song selection within a songlist.
type Select struct {
	songlistWidget func() *widgets.SonglistWidget
	toggle         bool
	finished       bool
}

func NewSelect(songlistWidget func() *widgets.SonglistWidget) *Select {
	return &Select{songlistWidget: songlistWidget}
}

func (cmd *Select) Reset() {
	cmd.finished = false
}

func (cmd *Select) Execute(t lexer.Token) error {
	var err error

	s := t.String()
	songlistWidget := cmd.songlistWidget()

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
		index := songlistWidget.Cursor()
		selected := songlistWidget.List().Selected(index)
		songlistWidget.List().SetSelected(index, !selected)
		songlistWidget.MoveCursor(1)

	default:
		return fmt.Errorf("Unexpected '%s', expected END", s)
	}

	return err
}
