// package db provides a shared object containing all of PMS' data.
package db

import (
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
)

// Instance holds state related to mutable data within PMS, such as the current
// state of MPD, any songlists, clipboards, options.
type Instance struct {
	// mpd state
	mpdStatus   pms_mpd.PlayerStatus
	currentSong *song.Song

	// song lists
	queue      *songlist.Queue
	library    *songlist.Library
	songlists  []songlist.Songlist
	clipboards map[string]songlist.Songlist
	options    *options.Options
}

func New() *Instance {
	return &Instance{
		clipboards: make(map[string]songlist.Songlist, 0),
	}
}

func (db *Instance) Queue() *songlist.Queue {
	return db.queue
}

func (db *Instance) SetQueue(queue *songlist.Queue) {
	db.queue = queue
}
