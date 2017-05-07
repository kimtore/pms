package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
)

// Quit exits the program.
type Redraw struct {
	api API
}

func NewRedraw(api API) Command {
	return &Redraw{
		api: api,
	}
}

func (cmd *Redraw) Execute(t lexer.Token) error {
	ui := cmd.api.UI()
	switch t.Class {
	case lexer.TokenEnd:
		ui.App.PostFunc(func() {
			ui.Refresh()
		})
		return nil
	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}
}
