package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/gdamore/tcell/views"
)

// Quit exits the program.
type Redraw struct {
	app *views.Application
}

func NewRedraw(app *views.Application) *Redraw {
	return &Redraw{app: app}
}

func (cmd *Redraw) Reset() {
}

func (cmd *Redraw) Execute(t lexer.Token) error {
	switch t.Class {
	case lexer.TokenEnd:
		cmd.app.PostFunc(func() {
			cmd.app.Refresh()
		})
		return nil
	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}
}
