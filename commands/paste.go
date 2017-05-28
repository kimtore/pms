package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Paste inserts songs from the clipboard.
type Paste struct {
	newcommand
	api      api.API
	position int
}

// NewPaste returns Paste.
func NewPaste(api api.API) Command {
	return &Paste{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Paste) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteVerbs(lit)

	// Expect either "before" or "after".
	switch tok {
	case lexer.TokenIdentifier:
		switch lit {
		case "before":
			cmd.position = 0
		case "after":
			cmd.position = 1
		default:
			return fmt.Errorf("Unexpected '%s', expected position", lit)
		}
		cmd.setTabCompleteEmpty()
		return cmd.ParseEnd()

	// Fall back to "after" if no arguments given.
	case lexer.TokenEnd:
		cmd.position = 1

	default:
		return fmt.Errorf("Unexpected '%s', expected position", lit)
	}

	return nil
}

// Exec implements Command.
func (cmd *Paste) Exec() error {
	list := cmd.api.Songlist()
	cursor := list.Cursor()
	clipboard := cmd.api.Clipboard()

	err := list.InsertList(clipboard, cursor+cmd.position)
	cmd.api.ListChanged()

	if err != nil {
		return err
	}

	cmd.api.Message("%d more tracks", clipboard.Len())

	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Paste) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"after",
		"before",
	})
}
