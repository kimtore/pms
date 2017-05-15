package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/constants"
	"github.com/ambientsound/pms/input/lexer"
)

// InputMode changes the Multibar's input mode.
type InputMode struct {
	command
	api  api.API
	mode int
}

func NewInputMode(api api.API) Command {
	return &InputMode{
		api: api,
	}
}

func (cmd *InputMode) Execute(class int, s string) error {
	multibar := cmd.api.Multibar()

	switch class {
	case lexer.TokenIdentifier:
		switch s {
		case "normal":
			cmd.mode = constants.MultibarModeNormal
		case "input":
			cmd.mode = constants.MultibarModeInput
		case "search":
			cmd.mode = constants.MultibarModeSearch
		default:
			cmd.mode = multibar.Mode()
		}
	case lexer.TokenEnd:
		multibar.SetMode(cmd.mode)

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return nil
}
