package topbar

import (
	"github.com/ambientsound/pms/api"
)

// Artist draws the short name of this application, as defined in the version module.
type Artist struct {
	tag string
	fragment
}

func NewArtist(a api.API) Fragment {
	return &Artist{
		"artist",
		fragment{api: a},
	}
}

func (w *Artist) Text() string {
	song := w.api.Song()
	if song == nil {
		return ""
	}
	text, ok := song.StringTags[w.tag]
	if !ok {
		text = "<unknown>"
	}
	return text
}

func (w *Artist) Width() int {
	return len(w.Text())
}

func (w *Artist) Draw(x, y int) int {
	return w.drawNextString(x, y, w.Text(), w.Style(w.tag))
}
