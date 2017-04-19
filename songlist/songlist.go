package songlist

import (
	"bufio"
	"io"
	"sort"
	"strings"

	"github.com/fhs/gompd/mpd"

	"github.com/ambientsound/pms/song"
)

type Songlist struct {
	Name                string
	Songs               []*song.Song
	currentSortCriteria string
}

func New() (s *Songlist) {
	s = &Songlist{}
	s.Songs = make([]*song.Song, 0)
	return
}

func (s *Songlist) Add(song *song.Song) {
	s.Songs = append(s.Songs, song)
}

func (s *Songlist) Sort() {
	sort.Sort(s)
	sort.Stable(s)
}

func (s *Songlist) Len() int {
	return len(s.Songs)
}

func (s *Songlist) Less(a, b int) bool {
	return s.Songs[a].TagString(s.currentSortCriteria) < s.Songs[b].TagString(s.currentSortCriteria)
}

func (s *Songlist) Swap(a, b int) {
	c := s.Songs[a]
	s.Songs[a] = s.Songs[b]
	s.Songs[b] = c
}

func NewFromFile(file io.Reader) (songs *Songlist) {
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
			s.Tags[tokens[0]] = []rune(tokens[1])
		}
	}
	return
}

func NewFromAttrlist(attrlist []mpd.Attrs) *Songlist {
	songs := New()
	for _, attrs := range attrlist {
		s := song.New()
		s.SetTags(attrs)
		songs.Add(s)
	}
	return songs
}
