package songlist

import (
	"github.com/ambientsound/pms/song"
)

// Queue is a Songlist which represents the MPD play queue.
type Queue struct {
	BaseSonglist
}

func NewQueue() (s *Queue) {
	s = &Queue{}
	s.songs = make([]*song.Song, 0)
	return
}
