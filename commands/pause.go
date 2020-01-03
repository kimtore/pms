package commands

import (
	"github.com/ambientsound/pms/api"
)

// Pause toggles play/paused state.
type Pause struct {
	newcommand
	api api.API
}

func NewPause(api api.API) Command {
	return &Pause{
		api: api,
	}
}

func (cmd *Pause) Parse() error {
	return cmd.ParseEnd()
}

func (cmd *Pause) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	// FIXME: play if paused
	return client.Pause()
}
