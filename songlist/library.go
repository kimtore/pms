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
	return fmt.Errorf("The song library is cannot be cleared because it is read-only. For a more effective method, try 'rm -rf'")
}

func (s *Library) Sort(fields []string) error {
	return fmt.Errorf("The song library is read-only. Please make a copy if you want to sort.")
}

func IsLibrary(s Songlist) bool {
	switch s.(type) {
	case *Library:
		return true
	default:
		return false
	}
}