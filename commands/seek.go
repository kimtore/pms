package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
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

	playerStatus := cmd.api.PlayerStatus()

	_, lit, absolute, err := cmd.ParseInt()
	if err != nil {
		return err
	}

	if absolute {
		cmd.absolute = lit
	} else {
		cmd.absolute = int(playerStatus.Elapsed) + lit
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
