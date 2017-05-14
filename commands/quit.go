package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Quit exits the program.
type Quit struct {
	command
	api api.API
}

func NewQuit(api api.API) Command {
	return &Quit{
		api: api,
	}
}

func (cmd *Quit) Execute(class int, s string) error {
	switch class {
	case lexer.TokenEnd:
		cmd.api.Quit()
		return nil
	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}
}
