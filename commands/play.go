package commands

import (
	"fmt"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Play plays songs in the MPD playlist.
type Play struct {
	newcommand
	api       api.API
	cursor    bool
	selection bool
}

// NewPlay returns Play.
func NewPlay(api api.API) Command {
	return &Play{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Play) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	cmd.setTabCompleteVerbs(lit)

	switch tok {
	case lexer.TokenEnd:
		// No parameters; just send 'play' command to MPD
		return nil
	case lexer.TokenIdentifier:
	default:
		return fmt.Errorf("Unexpected '%s', expected identifier", lit)
	}

	switch lit {
	// Play song under cursor
	case "cursor":
		cmd.cursor = true
	// Play selected songs
	case "selection":
		cmd.selection = true
	default:
		return fmt.Errorf("Unexpected '%s', expected identifier", lit)
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Play) Exec() error {

	// Ensure MPD connection.
	client := cmd.api.MpdClient()
	if client == nil {
		return fmt.Errorf("Cannot play: not connected to MPD")
	}

	switch {
	case cmd.cursor:
		// Play song under cursor.
		return cmd.playCursor(client)
	case cmd.selection:
		// Play selected songs.
		return cmd.playSelection(client)
	}

	// If a selection is not given, start playing with default parameters.
	return client.Play(-1)
}

// playCursor plays the song under the cursor.
func (cmd *Play) playCursor(client *mpd.Client) error {

	// Get the song under the cursor.
	song := cmd.api.Songlist().CursorSong()
	if song == nil {
		return fmt.Errorf("Cannot play: no song under cursor")
	}

	// Check if the currently selected song has an ID. If it doesn't, it's not
	// from the queue, and the song will have to be added beforehand.
	id := song.ID
	if song.NullID() {
		var err error
		id, err = client.AddID(song.StringTags["file"], -1)
		if err != nil {
			return err
		}
	}

	// Play the correct song.
	return client.PlayID(id)
}

// playSelection plays the currently selected songs.
func (cmd *Play) playSelection(client *mpd.Client) error {

	// Get the track selection.
	selection := cmd.api.Songlist().Selection()
	if selection.Len() == 0 {
		return fmt.Errorf("Cannot play: no selection")
	}

	// Check if the first song has an ID. If it does, just start playing. The
	// playback order cannot be guaranteed as the selection might be
	// fragmented, so don't touch the selection.
	first := selection.Song(0)
	if !first.NullID() {
		return client.PlayID(first.ID)
	}

	// We are not operating directly on the queue; add all songs to the queue now.
	queue := cmd.api.Queue()
	queueLen := queue.Len()
	err := queue.AddList(selection)
	if err != nil {
		return err
	}
	cmd.api.Songlist().ClearSelection()
	cmd.api.Message("Playing %d new songs", selection.Len())

	// We haven't got the ID from the first added song, so use positions
	// instead. In case of simultaneous operation with another client, this
	// might lead to a race condition. Ignore this for now.
	return client.Play(queueLen)
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Play) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"cursor",
		"selection",
	})
}
