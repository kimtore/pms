package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Select manipulates song selection within a songlist.
type Select struct {
	command
	api      api.API
	toggle   bool
	visual   bool
	finished bool
}

func NewSelect(api api.API) Command {
	return &Select{
		api: api,
	}
}

func (cmd *Select) Execute(class int, s string) error {
	var err error

	list := cmd.api.Songlist()

	switch class {

	case lexer.TokenIdentifier:
		if cmd.finished {
			return fmt.Errorf("Unexpected '%s', expected END", s)
		}
		switch s {
		// Toggle cursor select on/off
		case "toggle":
			cmd.toggle = true
		// Toggle visual mode on/off
		case "visual":
			cmd.visual = true
		default:
			return fmt.Errorf("Unexpected '%s', expected identifier", s)
		}
		cmd.finished = true

	case lexer.TokenEnd:
		if !cmd.finished {
			return fmt.Errorf("Unexpected END, expected identifier")
		}

		switch {
		case cmd.toggle && list.HasVisualSelection():
			list.CommitVisualSelection()
			list.DisableVisualSelection()

		case cmd.visual:
			list.ToggleVisualSelection()
			return nil

		default:
			index := list.Cursor()
			selected := list.Selected(index)
			list.SetSelected(index, !selected)
		}

		list.MoveCursor(1)

	default:
		return fmt.Errorf("Unexpected '%s', expected END", s)
	}

	return err
}
