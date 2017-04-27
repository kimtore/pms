package songlist

import (
	"fmt"
	"sort"
	"time"

	"github.com/ambientsound/gompd/mpd"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/song"
)

type Songlist interface {
	Add(*song.Song) error
	Clear() error
	Duplicate(Songlist) error
	Len() int
	Less(int, int) bool
	Name() string
	Replace(int, *song.Song) error
	SetName(string) error
	Song(int) *song.Song
	Songs() []*song.Song
	Sort([]string) error
	Swap(int, int)
	Truncate(int) error
}

type BaseSonglist struct {
	name                string
	songs               []*song.Song
	currentSortCriteria string
}

func New() (s *BaseSonglist) {
	s = &BaseSonglist{}
	s.Clear()
	return
}

func (s *BaseSonglist) Add(song *song.Song) error {
	s.songs = append(s.songs, song)
	return nil
}

func (s *BaseSonglist) Clear() error {
	s.songs = make([]*song.Song, 0)
	return nil
}

func (s *BaseSonglist) Replace(index int, song *song.Song) error {
	if index < 0 || index >= s.Len() {
		return fmt.Errorf("Out of bounds")
	}
	s.songs[index] = song
	return nil
}

// Duplicate makes a copy of the current songlist, and places it in dest.
func (s *BaseSonglist) Duplicate(dest Songlist) error {
	if err := dest.Clear(); err != nil {
		return err
	}
	oldSongs := s.Songs()
	for i := range oldSongs {
		song := *oldSongs[i]
		if err := dest.Add(&song); err != nil {
			return err
		}
	}
	return nil
}

func (s *BaseSonglist) Truncate(length int) error {
	if length < 0 || length > s.Len() {
		return fmt.Errorf("Out of bounds")
	}
	s.songs = s.songs[:length]
	return nil
}

func (s *BaseSonglist) SetName(name string) error {
	s.name = name
	return nil
}

func (s *BaseSonglist) Name() string {
	return s.name
}

func (s *BaseSonglist) Song(i int) *song.Song {
	if i < 0 || i >= s.Len() {
		return nil
	}
	return s.songs[i]
}

func (s *BaseSonglist) Songs() []*song.Song {
	return s.songs
}

func (s *BaseSonglist) Sort(fields []string) error {
	if len(fields) == 0 {
		return fmt.Errorf("Cannot sort without sort criteria")
	}
	s.sortBy(fields[0])
	for _, field := range fields[1:] {
		s.stableSortBy(field)
	}
	return nil
}

func (s *BaseSonglist) sortBy(field string) {
	s.currentSortCriteria = field
	timer := time.Now()
	sort.Sort(s)
	console.Log("Sorted '%s' by '%s' in %s", s.Name(), field, time.Since(timer).String())
}

func (s *BaseSonglist) stableSortBy(field string) {
	s.currentSortCriteria = field
	timer := time.Now()
	sort.Stable(s)
	console.Log("Stable sorted '%s' by '%s' in %s", s.Name(), field, time.Since(timer).String())
}

func (s *BaseSonglist) Len() int {
	return len(s.songs)
}

func (s *BaseSonglist) Less(a, b int) bool {
	return s.songs[a].SortTags[s.currentSortCriteria] < s.songs[b].SortTags[s.currentSortCriteria]
}

func (s *BaseSonglist) Swap(a, b int) {
	s.songs[a], s.songs[b] = s.songs[b], s.songs[a]
}

func (songs *BaseSonglist) AddFromAttrlist(attrlist []mpd.Attrs) {
	for _, attrs := range attrlist {
		s := song.New()
		s.SetTags(attrs)
		songs.Add(s)
	}
}
