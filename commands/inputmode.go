package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/widgets"
)

// InputMode changes the Multibar's input mode.
type InputMode struct {
	api  API
	mode int
}

func NewInputMode(api API) Command {
	return &InputMode{
		api: api,
	}
}

func (cmd *InputMode) Execute(t lexer.Token) error {
	s := t.String()
	multibar := cmd.api.Multibar()

	switch t.Class {
	case lexer.TokenIdentifier:
		switch s {
		case "normal":
			cmd.mode = widgets.MultibarModeNormal
		case "visual":
			switch multibar.Mode() {
			case widgets.MultibarModeVisual:
				cmd.mode = widgets.MultibarModeNormal
			default:
				cmd.mode = widgets.MultibarModeVisual
			}
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
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return nil
}
