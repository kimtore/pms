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
	AddList(Songlist) error
	Clear() error
	Delete() error
	Duplicate(Songlist) error
	InRange(int) bool
	Len() int
	Less(int, int) bool
	Locate(*song.Song) (int, error)
	Lock()
	Name() string
	Remove(int) error
	RemoveIndices([]int) error
	Replace(int, *song.Song) error
	SetName(string) error
	Song(int) *song.Song
	Songs() []*song.Song
	Sort([]string) error
	Swap(int, int)
	Truncate(int) error
	Unlock()

	ClearSelection()
	Columns([]string) Columns
	CommitVisualSelection()
	Cursor() int
	CursorSong() *song.Song
	CursorToSong(*song.Song) error
	DisableVisualSelection()
	EnableVisualSelection()
	HasVisualSelection() bool
	IndexAtSong(int, *song.Song) bool
	MoveCursor(int)
	Selected(int) bool
	Selection() Songlist
	SelectionIndices() []int
	SetCursor(int)
	SetSelected(int, bool)
	ToggleVisualSelection()
	ValidateCursor(int, int)
}

type BaseSonglist struct {
	name                string
	songs               []*song.Song
	currentSortCriteria string
	mutex               sync.Mutex

	columns         ColumnMap
	cursor          int
	selection       map[int]struct{}
	visualSelection [3]int
}

func New() (s *BaseSonglist) {
	s = &BaseSonglist{}
	s.clear()
	return
}

// Add adds a song to the songlist.
func (s *BaseSonglist) Add(song *song.Song) error {
	s.add(song)
	return nil
}

// add internally adds a song to the songlist, without any side effects at MPD's side.
func (s *BaseSonglist) add(song *song.Song) {
	s.songs = append(s.songs, song)
	s.ensureColumns(song)
	s.columns.Add(song)
}

// AddList appends a songlist to this songlist.
func (s *BaseSonglist) AddList(songlist Songlist) error {
	songs := songlist.Songs()
	for _, song := range songs {
		if err := s.Add(song); err != nil {
			return err
		}
	}
	return nil
}

func (s *BaseSonglist) Remove(index int) error {
	song := s.Song(index)
	if song == nil {
		return fmt.Errorf("Out of bounds")
	}

	s.columns.Remove(song)
	s.Lock()
	defer s.Unlock()

	//console.Log("Removing song number %d from songlist '%s'", index, s.Name())
	if index+1 == s.Len() {
		s.songs = s.songs[:index]
	} else {
		s.songs = append(s.songs[:index], s.songs[index+1:]...)
	}
	return nil
}

// RemoveIndices removes a selection of songs from the songlist, having the
// index defined by the int slice parameter.
func (s *BaseSonglist) RemoveIndices(indices []int) error {
	// Ensure that indices are removed in reverse order
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))
	for _, i := range indices {
		if err := s.Remove(i); err != nil {
			return err
		}
	}
	return nil
}

func (s *BaseSonglist) Lock() {
	s.mutex.Lock()
}

func (s *BaseSonglist) Unlock() {
	s.mutex.Unlock()
}

// Clear clears the songlist by removing any songs.
func (s *BaseSonglist) Clear() error {
	s.clear()
	return nil
}

func (s *BaseSonglist) clear() {
	s.songs = make([]*song.Song, 0)
	s.columns = make(ColumnMap)
	s.ClearSelection()
}

// Delete deletes a songlist. This is a placeholder function that should be
// overridden by other classes that need to trigger an action on the MPD side.
func (s *BaseSonglist) Delete() error {
	return nil
}

func (s *BaseSonglist) Replace(index int, song *song.Song) error {
	if !s.InRange(index) {
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
	if match == nil {
		return 0, fmt.Errorf("Attempt to locate nil song")
	}
	for i, test := range s.songs {
		hasId := !(match.NullID() || test.NullID())
		switch {
		case hasId && match.ID == test.ID:
		case !hasId && match.StringTags["file"] == test.StringTags["file"]:
		default:
			continue
		}
		return i, nil
	}
	return 0, fmt.Errorf("Cannot find song in songlist '%s'", s.Name())
}

// IndexAtSong returns true if the song at the specified index is at a song
// with the same ID, or the path if the songlist is not a queue.
func (s *BaseSonglist) IndexAtSong(i int, song *song.Song) bool {
	check := s.Song(i)
	return song != nil && check != nil && check.StringTags["file"] == song.StringTags["file"]
}

func (s *BaseSonglist) SetName(name string) error {
	s.name = name
	return nil
}

func (s *BaseSonglist) Name() string {
	return s.name
}

func (s *BaseSonglist) Song(i int) *song.Song {
	if !s.InRange(i) {
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

// InRange returns true if the provided index is within songlist range, false otherwise.
func (s *BaseSonglist) InRange(index int) bool {
	return index >= 0 && index < s.Len()
}

func (s *BaseSonglist) Less(a, b int) bool {
	return s.songs[a].SortTags[s.currentSortCriteria] < s.songs[b].SortTags[s.currentSortCriteria]
}

func (s *BaseSonglist) Swap(a, b int) {
	s.songs[a], s.songs[b] = s.songs[b], s.songs[a]
}

func (s *BaseSonglist) AddFromAttrlist(attrlist []mpd.Attrs) {
	for _, attrs := range attrlist {
		newSong := song.New()
		newSong.SetTags(attrs)
		s.add(newSong)
	}
}
