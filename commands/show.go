package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Show directs which window (main widget) to show.
type Show struct {
	newcommand
	api    api.API
	window api.Window
}

// NewShow returns Show.
func NewShow(api api.API) Command {
	return &Show{
		api: api,
	}
}

// Parse parses the viewport movement command.
func (cmd *Show) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteVerbs(lit)

	switch tok {
	case lexer.TokenIdentifier:
	default:
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	switch lit {
	case "logs":
		cmd.window = api.WindowLogs
	case "music":
		cmd.window = api.WindowMusic
	case "playlists":
		cmd.window = api.WindowPlaylists
	default:
		return fmt.Errorf("can't show '%s'; no such window", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Show) Exec() error {
	cmd.api.UI().ActivateWindow(cmd.window)
	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Show) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"logs",
		"music",
		"playlists",
	})
}
