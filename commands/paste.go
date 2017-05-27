package commands

import (
	"github.com/ambientsound/pms/api"
)

// Paste inserts songs from the clipboard.
type Paste struct {
	newcommand
	api api.API
}

// NewPaste returns Paste.
func NewPaste(api api.API) Command {
	return &Paste{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Paste) Parse() error {
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Paste) Exec() error {
	list := cmd.api.Songlist()
	cursor := list.Cursor()
	clipboard := cmd.api.Clipboard()

	err := list.InsertList(clipboard, cursor+1)
	cmd.api.ListChanged()

	if err != nil {
		cmd.api.Message("%d more tracks", clipboard.Len())
	}

	return err
}
