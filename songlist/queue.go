package songlist

import (
	"fmt"

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

func (s *Queue) Name() string {
	return "Queue"
}

func (s *Queue) SetName(name string) error {
	return fmt.Errorf("The queue name cannot be changed.")
}

func (s *Queue) Clear() error {
	return fmt.Errorf("Clearing the queue is not implemented yet.")
}

func (s *Queue) Sort(fields []string) error {
	return fmt.Errorf("Sorting the queue is not implemented yet.")
}

// Merge incorporates songs from another songlist, replacing songs that has the same position.
func (q *Queue) Merge(s Songlist) (*Queue, error) {
	newQueue := NewQueue()

	oldSongs := q.Songs()
	for i := range oldSongs {
		song := *oldSongs[i]
		newQueue.Add(&song)
	}

	newSongs := s.Songs()
	for i := range newSongs {
		song := newSongs[i]
		switch {
		case song.Position < 0:
			return nil, fmt.Errorf("Song number %d does not have a position", i)
		case song.Position == newQueue.Len():
			if err := newQueue.Add(song); err != nil {
				return nil, fmt.Errorf("Cannot add song %d to queue: %s", i, err)
			}
		case song.Position > newQueue.Len():
			return nil, fmt.Errorf("Song number %d has position greater than list length, there are parts missing", i)
		default:
			if err := newQueue.Replace(song.Position, song); err != nil {
				return nil, fmt.Errorf("Cannot replace song %d: %s", i, err)
			}
		}
	}

	return newQueue, nil
}

func IsQueue(s Songlist) bool {
	switch s.(type) {
	case *Queue:
		return true
	default:
		return false
	}
}
