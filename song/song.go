package song

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ambientsound/pms/utils"

	"github.com/ambientsound/gompd/mpd"
)

// Song represents a combined view of a song from both MPD and PMS' perspectives.
type Song struct {
	ID       int
	Position int
	Time     int
	Tags     Taglist
	SortTags StringTaglist
}

type Tag []rune

type Taglist map[string]Tag

type StringTaglist map[string]string

func New() (s *Song) {
	s = &Song{}
	s.Tags = make(Taglist)
	s.SortTags = make(StringTaglist)
	return
}

func (s *Song) TagString(key string) string {
	return string(s.Tags[key])
}

func (s *Song) SetTags(tags mpd.Attrs) {
	s.Tags = make(Taglist)
	for key := range tags {
		lowKey := strings.ToLower(key)
		s.Tags[lowKey] = []rune(tags[key])
	}
	s.AutoFill()
	s.FillSortTags()
}

// AutoFill post-processes and caches song tags.
func (s *Song) AutoFill() {
	var err error

	s.ID, _ = strconv.Atoi(s.TagString("id"))
	s.Position, _ = strconv.Atoi(s.TagString("pos"))

	s.Time, err = strconv.Atoi(s.TagString("time"))
	if err == nil {
		s.Tags["time"] = utils.TimeRunes(s.Time)
	} else {
		s.Tags["time"] = utils.TimeRunes(-1)
	}
	if len(s.Tags["date"]) >= 4 {
		s.Tags["year"] = s.Tags["date"][:4]
	}
}

// FillSortTags post-processes tags, and saves them as strings for sorting purposes later on.
func (s *Song) FillSortTags() {
	for i := range s.Tags {
		s.SortTags[i] = strings.ToLower(s.TagString(i))
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
