package song

import (
	"strconv"

	"github.com/ambientsound/pms/utils"

	"github.com/fhs/gompd/mpd"
)

type Song struct {
	ID       int
	Position int
	Time     int
	Tags     mpd.Attrs
}

func New() (s *Song) {
	s = &Song{}
	s.Tags = make(mpd.Attrs)
	return
}

func (s *Song) SetTags(tags mpd.Attrs) {
	s.Tags = tags
	s.AutoFill()
}

func (s *Song) AutoFill() {
	var err error
	s.Time, err = strconv.Atoi(s.Tags["Time"])
	if err == nil {
		s.Tags["Time"] = utils.TimeString(s.Time)
	} else {
		s.Tags["Time"] = `--:--`
	}
}
