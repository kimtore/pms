package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Select manipulates song selection within a songlist.
type Select struct {
	newcommand
	api    api.API
	toggle bool
	visual bool
	nearby []string
}

// NewSelect returns Select.
func NewSelect(api api.API) Command {
	return &Select{
		api:    api,
		nearby: make([]string, 0),
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
	case "toggle":
		cmd.toggle = true
	case "visual":
		cmd.visual = true
	case "nearby":
		return cmd.parseNearby()
	default:
		return fmt.Errorf("Unexpected '%s', expected identifier", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Select) Exec() error {
	list := cmd.api.UI().TableWidget().List()

	switch {
	case cmd.toggle && list.HasVisualSelection():
		list.CommitVisualSelection()
		list.DisableVisualSelection()

	case cmd.visual:
		list.ToggleVisualSelection()
		return nil

	case len(cmd.nearby) > 0:
		return cmd.selectNearby()

	default:
		index := list.Cursor()
		selected := list.Selected(index)
		list.SetSelected(index, !selected)
	}

	list.MoveCursor(1)

	return nil
}

// parseNearby parses tags and inserts them in the nearby list.
func (cmd *Select) parseNearby() error {

	// Data initialization and sanity checks
	list := cmd.api.Songlist()
	song := list.CursorSong()

	// Retrieve a list of songs
	tags, err := cmd.ParseTags(song.TagKeys())
	if err != nil {
		return err
	}

	cmd.nearby = tags
	return nil
}

// selectNearby selects tracks near the cursor with similar tags.
func (cmd *Select) selectNearby() error {
	list := cmd.api.Songlist()
	index := list.Cursor()
	song := list.CursorSong()

	// In case the list has a visual selection, disable that selection instead.
	if list.HasVisualSelection() {
		list.DisableVisualSelection()
		return nil
	}

	if song == nil {
		return fmt.Errorf("Can't select nearby songs; no song under cursor")
	}

	// Find the start and end positions
	start := list.NextOf(cmd.nearby, index+1, -1)
	end := list.NextOf(cmd.nearby, index, 1) - 1

	// Set visual selection and move cursor to end of selection
	list.SetVisualSelection(start, end, start)
	list.SetCursor(end)

	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Select) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"nearby",
		"toggle",
		"visual",
	})
}
