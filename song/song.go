package song

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ambientsound/pms/utils"

	"github.com/ambientsound/gompd/mpd"
)

// Song represents a combined view of a song from both MPD and PMS' perspectives.
type Song struct {
	ID         int
	Position   int
	Time       int
	Tags       Taglist
	StringTags StringTaglist
	SortTags   StringTaglist
}

type Tag []rune

type Taglist map[string]Tag

type StringTaglist map[string]string

const NullID int = -1
const NullPosition int = -1

func New() (s *Song) {
	s = &Song{}
	s.Tags = make(Taglist)
	s.StringTags = make(StringTaglist)
	s.SortTags = make(StringTaglist)
	return
}

func (s *Song) SetTags(tags mpd.Attrs) {
	s.Tags = make(Taglist)
	for key := range tags {
		lowKey := strings.ToLower(key)
		s.Tags[lowKey] = []rune(tags[key])
		s.StringTags[lowKey] = tags[key]
	}
	s.AutoFill()
	s.FillSortTags()
}

// NullID returns true if the song's ID is not present.
func (s *Song) NullID() bool {
	return s.ID == NullID
}

// NullPosition returns true if the song's osition is not present.
func (s *Song) NullPosition() bool {
	return s.Position == NullPosition
}

// AutoFill post-processes and caches song tags.
func (s *Song) AutoFill() {
	var err error

	s.ID, err = strconv.Atoi(s.StringTags["id"])
	if err != nil {
		s.ID = NullID
	}
	s.Position, err = strconv.Atoi(s.StringTags["pos"])
	if err != nil {
		s.Position = NullPosition
	}

	s.Time, err = strconv.Atoi(s.StringTags["time"])
	if err == nil {
		s.Tags["time"] = utils.TimeRunes(s.Time)
	} else {
		s.Tags["time"] = utils.TimeRunes(-1)
	}
	if len(s.Tags["date"]) >= 4 {
		s.Tags["year"] = s.Tags["date"][:4]
		s.StringTags["year"] = string(s.Tags["year"])
	}
}

// FillSortTags post-processes tags, and saves them as strings for sorting purposes later on.
func (s *Song) FillSortTags() {
	for i := range s.Tags {
		s.SortTags[i] = strings.ToLower(s.StringTags[i])
	}

	if t, ok := s.SortTags["track"]; ok {
		s.SortTags["track"] = trackSort(t)
	}

	if _, ok := s.SortTags["artistsort"]; !ok {
		s.SortTags["artistsort"] = s.SortTags["artist"]
	}

	if _, ok := s.SortTags["albumartist"]; !ok {
		s.SortTags["albumartist"] = s.SortTags["artist"]
	}

	if _, ok := s.SortTags["albumartistsort"]; !ok {
		s.SortTags["albumartistsort"] = s.SortTags["albumartist"]
	}
}

// HasOneOfTags returns true if the song contains at least one of the tags mentioned.
func (s *Song) HasOneOfTags(tags ...string) bool {
	for _, tag := range tags {
		if _, ok := s.Tags[tag]; ok {
			return true
		}
	}
	return false
}

// TagKeys returns a string slice with all tag keys, sorted in alphabetical order.
func (s *Song) TagKeys() []string {
	keys := make(sort.StringSlice, 0, len(s.StringTags))
	for tag := range s.StringTags {
		keys = append(keys, tag)
	}
	keys.Sort()
	return keys
}

func trackSort(s string) string {
	tracks := strings.Split(s, "/")
	if len(tracks) == 0 {
		return s
	}
	trackNum, err := strconv.Atoi(tracks[0])
	if err != nil {
		return s
	}
	// Assume no release has more than 999 tracks.
	return fmt.Sprintf("%03d", trackNum)
}
