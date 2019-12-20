package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/constants"
	"github.com/ambientsound/pms/input/lexer"
)

// InputMode changes the Multibar's input mode.
type InputMode struct {
	newcommand
	api  api.API
	mode constants.InputMode
}

func NewInputMode(api api.API) Command {
	return &InputMode{
		api: api,
	}
}

func (cmd *InputMode) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	switch tok {
	case lexer.TokenIdentifier:
		break
	default:
		return fmt.Errorf("unexpected '%s'; expected identifier", lit)
	}

	switch lit {
	case "normal":
		cmd.mode = constants.MultibarModeNormal
	case "input":
		cmd.mode = constants.MultibarModeInput
	case "search":
		cmd.mode = constants.MultibarModeSearch
	default:
		return fmt.Errorf("invalid input mode '%s'; expected one of 'normal', 'input', 'search'")
	}

	err := cmd.ParseEnd()
	if err != nil {
		return err
	}

	return nil
}

func (cmd *InputMode) Exec() error {
	cmd.api.SetInputMode(cmd.mode)
	return nil
}
