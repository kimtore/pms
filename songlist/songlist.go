package songlist

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/fhs/gompd/mpd"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/song"
)

type Songlist interface {
	Add(*song.Song) error
	AddList(Songlist) error
	Clear() error
	Delete() error
	Duplicate(Songlist) error
	Indices([]int) Songlist
	InRange(int) bool
	Insert(*song.Song, int) error
	InsertList(Songlist, int) error
	Len() int
	Locate(*song.Song) (int, error)
	Lock()
	Name() string
	NextOf([]string, int, int) int
	Remove(int) error
	RemoveIndices([]int) error
	Replace(int, *song.Song) error
	SetName(string) error
	Song(int) *song.Song
	Songs() []*song.Song
	Sort([]string) error
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
	SetUpdated()
	SetVisualSelection(int, int, int)
	ToggleVisualSelection()
	Updated() time.Time
	ValidateCursor(int, int)
}

type BaseSonglist struct {
	name    string
	songs   []*song.Song
	mutex   sync.Mutex
	updated time.Time

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

// insert internally inserts a song to the songlist, without any side effects at MPD's side.
func (s *BaseSonglist) insert(newSong *song.Song, position int) {
	// create a copy of the list, with the new song inserted at the correct position
	songs := make([]*song.Song, len(s.songs)+1)
	position += copy(songs[:], s.songs[:position])
	songs[position] = newSong
	copy(songs[position+1:], s.songs[position:])
	s.songs = songs

	// expand columns
	s.ensureColumns(newSong)
	s.columns.Add(newSong)
}

// Insert inserts a song at the specified position.
func (s *BaseSonglist) Insert(song *song.Song, position int) error {
	s.insert(song, position)
	if position <= s.Cursor() {
		s.MoveCursor(1)
	}
	return nil
}

// InsertList inserts the songs in a songlist into this songlist, at a specified position.
func (s *BaseSonglist) InsertList(list Songlist, position int) error {
	size := list.Len()
	//console.Log("making size %d array", s.Len()+size)
	songs := make([]*song.Song, s.Len()+size)
	//console.Log("copy to 0 <- songs[:%d]", position)
	copy(songs[:], s.songs[:position])
	//console.Log("copy to %d <- list", position)
	copy(songs[position:], list.Songs())
	//console.Log("copy to %d <- songs[%d:]", position+size, position)
	copy(songs[position+size:], s.songs[position:])
	s.songs = songs

	// Move cursor if pasting above.
	if position <= s.Cursor() {
		s.MoveCursor(size)
	}

	// expand columns
	for _, newSong := range list.Songs() {
		s.ensureColumns(newSong)
		s.columns.Add(newSong)
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
		newSong := *oldSongs[i]
		newSong.ID = song.NullID
		newSong.Position = song.NullPosition
		if err := dest.Add(&newSong); err != nil {
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

// Indices accepts a slice of integers pointing to song positions, returns a
// new songlist with those songs. Any mismatching integers are ignored.
func (s *BaseSonglist) Indices(indices []int) Songlist {
	dest := New()
	for _, i := range indices {
		if song := s.Song(i); song != nil {
			dest.Add(song)
		} else {
			console.Log("BUG: Indices() returned an integer '%d' that resulted in a nil song, ignoring", i)
		}
	}
	return dest
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

// NextOf searches forwards or backwards for songs having different tags than the specified song.
// The index of the next song is returned.
func (s *BaseSonglist) NextOf(tags []string, index int, direction int) int {
	offset := func(i int) int {
		if direction > 0 || i == 0 {
			return 0
		}
		return 1
	}

	len := s.Len()
	index -= offset(index)
	check := s.Song(index)
	//console.Log("Starting at index %d, cursor at %d, direction %d", index, s.Cursor(), direction)

	//console.Log("Check: %+v", check.StringTags)
	for ; index < len && index >= 0; index += direction {
		//console.Log("trying %d", index)
		song := s.Song(index)
		//console.Log(".....: %+v", song.StringTags)
		if song == nil {
			//console.Log("NextOf: empty song, break")
			break
		}
		for _, tag := range tags {
			if check.StringTags[tag] != song.StringTags[tag] {
				//console.Log("NextOf: tag '%s' on source '%s' differs from destination '%s', breaking", tag, check.StringTags[tag], song.StringTags[tag])
				return index + offset(index)
			}
		}
	}

	//console.Log("NextOf: fallthrough")
	return index + offset(index)
}

// Sort sorts the songlist by the given tags. The first tag is sorted normally,
// while the remaining tags are used for stable sorting.
func (s *BaseSonglist) Sort(fields []string) error {
	if len(fields) == 0 {
		return fmt.Errorf("Cannot sort without sort criteria")
	}
	s.Lock()
	defer s.Unlock()

	stable := false
	for _, field := range fields {
		s.sortBy(field, stable)
		stable = true
	}

	return nil
}

// sortBy sorts the songlist by the given tag, optionally using stable sort.
func (s *BaseSonglist) sortBy(field string, stable bool) {
	sortFunc := func(a, b int) bool {
		return s.songs[a].SortTags[field] < s.songs[b].SortTags[field]
	}
	timer := time.Now()
	if stable {
		sort.SliceStable(s.songs, sortFunc)
	} else {
		sort.Slice(s.songs, sortFunc)
	}
	console.Log("Sorted '%s' by '%s' in %s", s.Name(), field, time.Since(timer).String())
}

func (s *BaseSonglist) Len() int {
	return len(s.songs)
}

// InRange returns true if the provided index is within songlist range, false otherwise.
func (s *BaseSonglist) InRange(index int) bool {
	return index >= 0 && index < s.Len()
}

func (s *BaseSonglist) AddFromAttrlist(attrlist []mpd.Attrs) {
	for _, attrs := range attrlist {
		newSong := song.New()
		newSong.SetTags(attrs)
		s.add(newSong)
	}
}

// Updated returns the timestamp of when this songlist was last updated.
func (s *BaseSonglist) Updated() time.Time {
	return s.updated
}

// SetUpdated sets the update timestamp of the songlist.
func (s *BaseSonglist) SetUpdated() {
	s.updated = time.Now()
}
