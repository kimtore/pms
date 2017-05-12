package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/widgets"
)

// InputMode changes the Multibar's input mode.
type InputMode struct {
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
			cmd.mode = widgets.MultibarModeNormal
		case "input":
			cmd.mode = widgets.MultibarModeInput
		case "search":
			cmd.mode = widgets.MultibarModeSearch
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
