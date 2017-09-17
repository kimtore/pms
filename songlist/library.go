package songlist

import (
	"fmt"
	"time"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/index"
)

// Library is a Songlist which represents the MPD song library.
type Library struct {
	BaseSonglist
	index   *index.Index
	version int
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

// OpenIndex configures the library to use the Bleve search index at the specified path.
func (s *Library) OpenIndex(path string) error {
	var err error

	if s.HasIndex() {
		if err = s.index.Close(); err != nil {
			return err
		}
		s.index = nil
	}

	s.index, err = index.New(path)

	return err
}

// CloseIndex closes the Bleve search index.
func (s *Library) CloseIndex() error {
	if s.HasIndex() {
		return s.index.Close()
	}
	return nil
}

func (s *Library) HasIndex() bool {
	return s.index != nil
}

func (s *Library) SetVersion(version int) {
	s.version = version
}

func (s *Library) Version() int {
	return s.version
}

func (s *Library) IndexSynced() bool {
	return s.HasIndex() && s.index.Version() == s.version
}

// FIXME: ReIndex is not thread safe yet!!!
func (s *Library) ReIndex() {
	go func() {
		timer := time.Now()
		s.index.IndexFull(s.Songs())
		s.index.SetVersion(s.Version())
		console.Log("Song library index complete, took %s", time.Since(timer).String())
	}()
}

// Search does a search in the Bleve index for a specific natural language
// query string, and returns a new Songlist with the search results.
func (s *Library) Search(q string) (Songlist, error) {
	if s.index == nil {
		return nil, fmt.Errorf("Search index is not open.")
	}

	ids, err := s.index.Search(q, s.Len())
	if err != nil {
		return nil, err
	}

	list := New()
	list.SetName(q)
	for _, id := range ids {
		song := s.Song(id)
		if song == nil {
			return nil, fmt.Errorf("Search index is corrupt.")
		}
		list.Add(song)
	}

	return list, nil
}

func (s *Library) Isolate(list Songlist, tags []string) (Songlist, error) {
	//names := make([]string, 0)
	//for k := range terms {
	//names = append(names, k)
	//}
	//name := strings.Join(names, ", ")
	//r.SetName(name)
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}
