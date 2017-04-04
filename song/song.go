package song

import (
	"github.com/fhs/gompd/mpd"
)

type Song struct {
	Tags mpd.Attrs
}

// Song type used to represent a song as it is found on the MPD server.
// FIXME
type MpdSong struct {
	ID       int
	Position int
	Song
}

func New() (s *Song) {
	s = &Song{}
	s.Tags = make(mpd.Attrs)
	return
}
