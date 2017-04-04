package songlist

import (
	"bufio"
	"io"
	"sort"
	"strings"

	"github.com/fhs/gompd/mpd"

	"github.com/ambientsound/pms/song"
)

type SongList struct {
	Name                string
	Songs               []*song.Song
	currentSortCriteria string
}

func New() (s *SongList) {
	s = &SongList{}
	s.Songs = make([]*song.Song, 0)
	return
}

func (s *SongList) Add(song *song.Song) {
	s.Songs = append(s.Songs, song)
}

func (s *SongList) Sort() {
	sort.Sort(s)
	sort.Stable(s)
}

func (s *SongList) Len() int {
	return len(s.Songs)
}

func (s *SongList) Less(a, b int) bool {
	return s.Songs[a].Tags[s.currentSortCriteria] < s.Songs[b].Tags[s.currentSortCriteria]
}

func (s *SongList) Swap(a, b int) {
	c := s.Songs[a]
	s.Songs[a] = s.Songs[b]
	s.Songs[b] = c
}

func NewFromFile(file io.Reader) (songs *SongList) {
	scanner := bufio.NewScanner(file)
	songs = New()
	var s *song.Song
	for scanner.Scan() {
		tokens := strings.SplitN(scanner.Text(), ": ", 2)
		if tokens[0] == "file" {
			if s != nil {
				songs.Add(s)
			}
			s = song.New()
		}
		if s != nil {
			s.Tags[tokens[0]] = tokens[1]
		}
	}
	return
}

func NewFromAttrlist(attrlist []mpd.Attrs) *SongList {
	songs := New()
	for _, attrs := range attrlist {
		s := &song.Song{}
		s.Tags = attrs
		songs.Add(s)
	}
	return songs
}
