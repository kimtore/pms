package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Single toggles MPD's single mode on and off.
type Single struct {
	newcommand
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
		return fmt.Errorf("Unexpected '%v', expected identifier", lit)
	}

	switch lit {
	case "on", "off", "toggle":
		break
	default:
		return fmt.Errorf("Unexpected '%v', expected identifier", lit)
	}

	cmd.action = lit

	cmd.setTabCompleteEmpty()
	return cmd.ParseEnd()

}

// Exec implements Command.
func (cmd *Single) Exec() error {

	client := cmd.api.MpdClient()
	if client == nil {
		return fmt.Errorf("Cannot change single mode: not connected to MPD.")
	}

	playerStatus := cmd.api.PlayerStatus()

	switch cmd.action {
	case "on":
		return client.Single(true)
	case "off":
		return client.Single(false)
	case "toggle", "":
		return client.Single(!playerStatus.Single)
	}

	return nil
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
