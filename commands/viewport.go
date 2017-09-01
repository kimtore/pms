package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Viewport acts on the viewport, such as scrolling the current songlist.
type Viewport struct {
	newcommand
	api      api.API
	relative int
}

// NewViewport returns Viewport.
func NewViewport(api api.API) Command {
	return &Viewport{
		api: api,
	}
}

// Parse parses the viewport movement command.
func (cmd *Viewport) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteVerbs(lit)

	switch tok {
	case lexer.TokenIdentifier:
	default:
		return fmt.Errorf("Unexpected '%s', expected identifier", lit)
	}

	switch lit {
	case "down":
		cmd.relative = 1
	case "up":
		cmd.relative = -1
	default:
		return fmt.Errorf("Viewport command '%s' not recognized", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Viewport) Exec() error {
	widget := cmd.api.SonglistWidget()

	switch {
	case cmd.relative != 0:
		widget.ScrollViewport(cmd.relative)
	}

	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Viewport) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"down",
		"up",
	})
}
