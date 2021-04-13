package commands

import (
	"fmt"

	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/zmb3/spotify"

	"github.com/ambientsound/pms/api"
)

// Add plays songs in the MPD playlist.
type Add struct {
	command
	api       api.API
	client    *spotify.Client
	tracklist *spotify_tracklist.List
}

// NewAdd returns Add.
func NewAdd(api api.API) Command {
	return &Add{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Add) Parse() error {
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Add) Exec() error {
	var err error

	cmd.tracklist = cmd.api.Tracklist()
	cmd.client, err = cmd.api.Spotify()

	if err != nil {
		return err
	}

	if cmd.tracklist == nil {
		return fmt.Errorf("cannot add to queue: not in a track list")
	}

	selection := cmd.tracklist.Selection()
	if selection.Len() == 0 {
		return fmt.Errorf("cannot add to queue: no selection")
	}

	// Allow command to deselect tracks in visual selection that were added to the queue.
	// In case of a queue add failure, it is desirable to still select the tracks that failed
	// to be added.
	cmd.tracklist.CommitVisualSelection()
	cmd.tracklist.DisableVisualSelection()

	for i, track := range selection.Tracks() {
		err := cmd.client.QueueSong(track.ID)
		if err != nil {
			return err
		}
		cmd.api.Message("Added %s to queue.", track.String())
		cmd.tracklist.SetSelected(i, false)
	}

	return nil
}
