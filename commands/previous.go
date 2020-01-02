package commands

import (
	"github.com/ambientsound/pms/api"
)

// Previous instructs the player to go to the previous song.
type Previous struct {
	newcommand
	api api.API
}

func NewPrevious(api api.API) Command {
	return &Previous{
		api: api,
	}
}

func (cmd *Previous) Parse() error {
	return cmd.ParseEnd()
}

func (cmd *Previous) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}
	return client.Previous()
}
