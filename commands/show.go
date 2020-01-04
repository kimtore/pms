package commands

import (
	"fmt"
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/log"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Show directs which window (main widget) to show.
type Show struct {
	command
	api  api.API
	list list.List
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
	case "selected":
		windows, ok := cmd.api.List().(*db.List)
		if !ok {
			return fmt.Errorf("`show selected` may only be used inside the windows view")
		}
		cmd.list = windows.Current()
	case "windows":
		cmd.list = cmd.api.Db()
	case "library":
		cmd.list = cmd.api.Library()
	case "logs":
		cmd.list = log.List(log.InfoLevel)
	default:
		return fmt.Errorf("can't show '%s'; no such window", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Show) Exec() error {
	cmd.api.SetList(cmd.list)
	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Show) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"library",
		"logs",
		"selected",
		"windows",
	})
}
