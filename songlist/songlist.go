package songlist

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ambientsound/gompd/mpd"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/song"
)

type Songlist interface {
	Add(*song.Song) error
	Clear() error
	Delete() error
	Duplicate(Songlist) error
	Len() int
	Less(int, int) bool
	Locate(*song.Song) (int, error)
	Lock()
	Name() string
	Remove(int) error
	Replace(int, *song.Song) error
	SetName(string) error
	Song(int) *song.Song
	Songs() []*song.Song
	Sort([]string) error
	Swap(int, int)
	Truncate(int) error
	Unlock()
}

type BaseSonglist struct {
	name                string
	songs               []*song.Song
	currentSortCriteria string
	mutex               sync.Mutex
}

func New() (s *BaseSonglist) {
	s = &BaseSonglist{}
	s.Clear()
	return
}

func (s *BaseSonglist) Add(song *song.Song) error {
	s.Lock()
	defer s.Unlock()
	s.songs = append(s.songs, song)
	return nil
}

func (s *BaseSonglist) Remove(index int) error {
	if !s.inRange(index) {
		return fmt.Errorf("List index out of range")
	}

	s.Lock()
	defer s.Unlock()

	if index+1 == s.Len() {
		s.songs = s.songs[:index]
	} else {
		s.songs = append(s.songs[:index], s.songs[index+1:]...)
	}
	return nil
}

func (s *BaseSonglist) Lock() {
	s.mutex.Lock()
}

func (s *BaseSonglist) Unlock() {
	s.mutex.Unlock()
}

func (s *BaseSonglist) Clear() error {
	s.songs = make([]*song.Song, 0)
	return nil
}

// Delete deletes a songlist. This is a placeholder function that should be
// overridden by other classes that need to trigger an action on the MPD side.
func (s *BaseSonglist) Delete() error {
	return nil
}

func (s *BaseSonglist) Replace(index int, song *song.Song) error {
	if !s.inRange(index) {
		return fmt.Errorf("Out of bounds")
	}
	s.Lock()
	defer s.Unlock()
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
	s.Lock()
	defer s.Unlock()
	s.songs = s.songs[:length]
	return nil
}

func (s *BaseSonglist) Locate(match *song.Song) (int, error) {
	for i, test := range s.songs {
		hasId := match.ID != -1 && test.ID != -1
		switch {
		case hasId && match.ID == test.ID:
		case match.StringTags["file"] == test.StringTags["file"]:
		default:
			continue
		}
		return i, nil
	}
	return 0, fmt.Errorf("Cannot find song in songlist %s", s.Name())
}

func (s *BaseSonglist) SetName(name string) error {
	s.name = name
	return nil
}

func (s *BaseSonglist) Name() string {
	return s.name
}

func (s *BaseSonglist) Song(i int) *song.Song {
	if !s.inRange(i) {
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
	s.Lock()
	defer s.Unlock()
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

func (s *BaseSonglist) inRange(index int) bool {
	return index >= 0 && index < s.Len()
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
