package commands

import (
	"fmt"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/multibar"
)

// InputMode changes the Multibar's input mode.
type InputMode struct {
	command
	api  api.API
	mode multibar.InputMode
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
		cmd.mode = multibar.ModeNormal
	case "input":
		cmd.mode = multibar.ModeInput
	case "search":
		cmd.mode = multibar.ModeSearch
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
	cmd.api.Multibar().SetMode(cmd.mode)
	return nil
}
