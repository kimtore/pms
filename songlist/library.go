package songlist

import (
	"github.com/ambientsound/pms/song"
)

// Library is a Songlist which represents the MPD song library.
type Library struct {
	BaseSonglist
}

func NewLibrary() (s *Library) {
	s = &Library{}
	s.songs = make([]*song.Song, 0)
	return
}

func (s *Library) Name() string {
	return "Library"
}

func IsLibrary(s Songlist) bool {
	switch s.(type) {
	case *Library:
		return true
	default:
		return false
	}
}
