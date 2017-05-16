package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Select manipulates song selection within a songlist.
type Select struct {
	newcommand
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

// Parse implements Command.
func (cmd *Select) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteVerbs(lit)

	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("Unexpected '%s', expected identifier", lit)
	}

	switch lit {
	// Toggle cursor select on/off
	case "toggle":
		cmd.toggle = true
	// Toggle visual mode on/off
	case "visual":
		cmd.visual = true
	default:
		return fmt.Errorf("Unexpected '%s', expected identifier", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

func (cmd *Select) Exec() error {
	list := cmd.api.Songlist()

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

	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Select) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"toggle",
		"visual",
	})
}
