package song

import (
	"strconv"
	"strings"

	"github.com/ambientsound/pms/utils"

	"github.com/fhs/gompd/mpd"
)

type Song struct {
	ID       int
	Position int
	Time     int
	Tags     Taglist
}

type Tag []rune

type Taglist map[string]Tag

func New() (s *Song) {
	s = &Song{}
	s.Tags = make(Taglist)
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
}

// AutoFill post-processes and caches song tags.
func (s *Song) AutoFill() {
	var err error
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
