package song

import (
	"fmt"
	"strconv"

	"github.com/fhs/gompd/mpd"
)

type Song struct {
	ID       int
	Position int
	Time     int
	Tags     mpd.Attrs
}

func TimeString(secs int) string {
	if secs < 0 {
		return "--:--"
	}
	hours := int(secs / 3600)
	secs = secs % 3600
	minutes := int(secs / 60)
	secs = secs % 60
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
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
		s.Tags["Time"] = TimeString(s.Time)
	} else {
		s.Tags["Time"] = `--:--`
	}
}
