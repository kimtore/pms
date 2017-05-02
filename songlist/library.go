package songlist

import (
	"fmt"

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

func (s *Library) SetName(name string) error {
	return fmt.Errorf("The song library name cannot be changed.")
}

func (s *Library) Clear() error {
	return fmt.Errorf("The song library cannot be cleared because it is read-only.")
}

func (s *Library) Delete() error {
	return fmt.Errorf("The song library cannot be deleted using PMS. Try 'rm -rf' in your favorite shell.")
}

func (s *Library) Sort(fields []string) error {
	return fmt.Errorf("The song library is read-only. Please make a copy if you want to sort.")
}

func (s *Library) Remove(index int) error {
	return fmt.Errorf("The song library is read-only.")
}

func (s *Library) RemoveIndices(indices []int) error {
	return fmt.Errorf("The song library is read-only.")
}

func IsLibrary(s Songlist) bool {
	switch s.(type) {
	case *Library:
		return true
	default:
		return false
	}
}
