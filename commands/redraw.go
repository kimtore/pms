package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Quit exits the program.
type Redraw struct {
	command
	api api.API
}

func NewRedraw(api api.API) Command {
	return &Redraw{
		api: api,
	}
}

func (cmd *Redraw) Execute(class int, s string) error {
	ui := cmd.api.UI()
	switch class {
	case lexer.TokenEnd:
		ui.PostFunc(func() {
			cmd.api.Db().Left().SetUpdated()
			cmd.api.Db().Right().SetUpdated()
			ui.Refresh()
		})
		return nil
	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}
}
