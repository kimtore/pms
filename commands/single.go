package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Single toggles MPD's single mode on and off.
type Single struct {
	command
	api    api.API
	action string
}

// NewSingle returns Single.
func NewSingle(api api.API) Command {
	return &Single{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Single) Parse() error {

	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteAction(lit)

	switch tok {
	case lexer.TokenIdentifier:
		break
	case lexer.TokenEnd:
		return nil
	default:
		return fmt.Errorf("unexpected '%v', expected identifier", lit)
	}

	playerStatus := cmd.api.PlayerStatus()

	switch lit {
	case "on":
		cmd.action = "single"
	case "off":
		cmd.action = "off"
	case "toggle":
		if playerStatus.RepeatState == "single" {
			cmd.action = "off"
		} else {
			cmd.action = "single"
		}
	default:
		return fmt.Errorf("unexpected '%v', expected identifier", lit)
	}

	cmd.action = lit

	cmd.setTabCompleteEmpty()
	return cmd.ParseEnd()

}

// Exec implements Command.
func (cmd *Single) Exec() error {

	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	return client.Repeat(cmd.action)
}

// setTabCompleteAction sets the tab complete list to available actions.
func (cmd *Single) setTabCompleteAction(lit string) {
	list := []string{
		"on",
		"off",
		"toggle",
	}
	cmd.setTabComplete(lit, list)
}
