package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Quit exits the program.
type Redraw struct {
	api api.API
}

func NewRedraw(api api.API) Command {
	return &Redraw{
		api: api,
	}
}

func (cmd *Redraw) Execute(t lexer.Token) error {
	ui := cmd.api.UI()
	switch t.Class {
	case lexer.TokenEnd:
		ui.PostFunc(func() {
			ui.Refresh()
		})
		return nil
	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}
}
