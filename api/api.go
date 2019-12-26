// Package api provides data model interfaces.
package api

import (
	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/input/keys"
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/ambientsound/pms/style"
	"github.com/zmb3/spotify"
)

// API defines a set of commands that should be available to commands run
// through the command-line interface.
type API interface {
	// Authenticate starts OAuth authentication.
	Authenticate()

	// Db returns the PMS database.
	Db() *db.Instance

	// Exec executes a command through the command-line interface.
	Exec(string) error

	// Return the global multibar instance.
	Multibar() *multibar.Multibar

	// Library returns the current MPD library, or nil if it has not been retrieved yet.
	Library() *songlist.Library

	// ListChanged notifies the UI that the current songlist has changed.
	ListChanged()

	// OptionChanged notifies that an option has been changed.
	OptionChanged(string)

	// Message sends a message to the user through the statusbar.
	Message(string, ...interface{})

	// MpdClient returns the current MPD client, which is confirmed to be alive. If the MPD connection is not working, nil is returned.
	MpdClient() *mpd.Client

	// Options returns PMS' global options.
	Options() *options.Options

	// PlayerStatus returns the current MPD player status.
	PlayerStatus() pms_mpd.PlayerStatus

	// Queue returns MPD's song queue.
	Queue() *songlist.Queue

	// Quit shuts down PMS.
	Quit()

	// Sequencer returns a pointer to the key sequencer that receives key events.
	Sequencer() *keys.Sequencer

	// Spotify returns a Spotify client.
	Spotify() spotify.Client

	// Song returns the currently playing song, or nil if no song is loaded.
	// Note that the song might be stopped, and the play/pause/stop status should
	// be checked using PlayerStatus().
	Song() *song.Song

	// Songlist returns the visible songlist.
	Songlist() songlist.Songlist

	// Styles returns the current stylesheet.
	Styles() style.Stylesheet

	// Tracklist returns the visible track list, if any.
	// Will be nil if the active widget shows a different kind of list.
	Tracklist() *spotify_tracklist.List

	// UI returns the global UI object.
	UI() UI
}
