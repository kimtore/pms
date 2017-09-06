package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Viewport acts on the viewport, such as scrolling the current songlist.
type Viewport struct {
	newcommand
	api        api.API
	movecursor bool
	relative   int
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
		cmd.movecursor = false
	case "up":
		cmd.relative = -1
		cmd.movecursor = false
	case "halfpgdn", "halfpagedn", "halfpagedown":
		cmd.scrollHalfPage(1)
	case "halfpgup", "halfpageup":
		cmd.scrollHalfPage(-1)
	case "pgdn", "pagedn", "pagedown":
		cmd.scrollFullPage(1)
	case "pgup", "pageup":
		cmd.scrollFullPage(-1)
	default:
		return fmt.Errorf("Viewport command '%s' not recognized", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// scrollHalfPage configures the command to scroll half a page up or down.
// The direction parameter must be -1 for up or 1 for down.
func (cmd *Viewport) scrollHalfPage(direction int) {
	_, y := cmd.api.SonglistWidget().Size()
	if y <= 1 {
		// Vim always moves at least one line
		cmd.relative = direction
	} else {
		cmd.relative = direction * y / 2
	}
	cmd.movecursor = true
}

// scrollFullPage configures the command to scroll a full page up or down.
// The direction parameter must be -1 for up or 1 for down.
func (cmd *Viewport) scrollFullPage(direction int) {
	_, y := cmd.api.SonglistWidget().Size()
	if y <= 3 {
		// Vim scrolls an entire page when 3 or fewer lines visible
		cmd.relative = direction * y
	} else if y == 4 {
		// Vim scrolls 3 lines when 4 lines visible
		cmd.relative = direction * 3
	} else {
		// Vim leaves 2 lines context when 5 or more lines visible
		cmd.relative = direction * (y - 2)
	}
	cmd.movecursor = false
}

// Exec implements Command.
func (cmd *Viewport) Exec() error {
	widget := cmd.api.SonglistWidget()

	widget.ScrollViewport(cmd.relative, cmd.movecursor)

	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Viewport) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"down",
		"halfpagedn",
		"halfpagedown",
		"halfpageup",
		"halfpgdn",
		"halfpgup",
		"pagedn",
		"pagedown",
		"pageup",
		"pgdn",
		"pgup",
		"up",
	})
}
