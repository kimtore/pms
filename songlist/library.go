package songlist

import (
	"fmt"
)

// Library is a Songlist which represents the MPD song library.
type Library struct {
	BaseSonglist
}

func NewLibrary() (s *Library) {
	s = &Library{}
	s.clear()
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

func (s *Library) Isolate(list Songlist, tags []string) (Songlist, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}
