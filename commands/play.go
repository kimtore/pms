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

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Play) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"cursor",
	})
}
