package commands

import (
	"github.com/ambientsound/pms/player"

	"github.com/ambientsound/pms/api"
)

// Seek seeks forwards or backwards in the currently playing track.
type Seek struct {
	command
	api          api.API
	absolute     int
	playerStatus player.State
}

// NewSeek returns Seek.
func NewSeek(api api.API) Command {
	return &Seek{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Seek) Parse() error {

	cmd.playerStatus = cmd.api.PlayerStatus()

	_, lit, absolute, err := cmd.ParseInt()
	if err != nil {
		return err
	}

	if absolute {
		cmd.absolute = lit * 1000
	} else {
		cmd.absolute = cmd.playerStatus.Progress + lit*1000
	}

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Seek) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	return client.Seek(cmd.absolute)
}
