package commands

import (
	"fmt"
	"strconv"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Seek seeks forwards or backwards in the currently playing track.
type Seek struct {
	newcommand
	api      api.API
	absolute int
}

// NewSeek returns Seek.
func NewSeek(api api.API) Command {
	return &Seek{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Seek) Parse() error {

	tok, lit := cmd.ScanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("Unexpected '%s', expected identifier", lit)
	}

	sign := 0
	playerStatus := cmd.api.PlayerStatus()

	// Relative seek
	switch lit[0] {
	case '+':
		sign = 1
		lit = lit[1:]
	case '-':
		sign = -1
		lit = lit[1:]
	}

	i, err := strconv.Atoi(lit)
	if err != nil {
		return fmt.Errorf("Unexpected '%s', expected number", lit)
	}

	if sign == 0 {
		cmd.absolute = i
	} else {
		cmd.absolute = int(playerStatus.Elapsed) + i*sign
	}

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Seek) Exec() error {
	mpdClient := cmd.api.MpdClient()
	if mpdClient == nil {
		return fmt.Errorf("Unable to set volume: cannot communicate with MPD")
	}

	playerStatus := cmd.api.PlayerStatus()
	return mpdClient.Seek(playerStatus.Song, cmd.absolute)
}
