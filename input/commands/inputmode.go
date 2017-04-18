package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/widgets"
)

// InputMode changes the Multibar's input mode.
type InputMode struct {
	multibarWidget *widgets.MultibarWidget
	mode           int
}

func NewInputMode(multibarWidget *widgets.MultibarWidget) *InputMode {
	return &InputMode{multibarWidget: multibarWidget}
}

func (cmd *InputMode) Reset() {
}

func (cmd *InputMode) Execute(t lexer.Token) error {
	s := string(t.Runes)

	switch t.Class {
	case lexer.TokenIdentifier:
		switch s {
		case "normal":
			cmd.mode = widgets.MultibarModeNormal
		case "input":
			cmd.mode = widgets.MultibarModeInput
		case "search":
			cmd.mode = widgets.MultibarModeSearch
		}
	case lexer.TokenEnd:
		cmd.multibarWidget.SetMode(cmd.mode)

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return nil
}
