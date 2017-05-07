package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
)

// Quit exits the program.
type Quit struct {
	api API
}

func NewQuit(api API) Command {
	return &Quit{
		api: api,
	}
}

func (cmd *Quit) Execute(t lexer.Token) error {
	switch t.Class {
	case lexer.TokenEnd:
		cmd.api.Quit()
		return nil
	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}
}
