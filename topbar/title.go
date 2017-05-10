package topbar

import (
	"github.com/ambientsound/pms/api"
)

// Title draws the current songlist's title.
type Title struct {
	api api.API
}

// NewTitle returns Title.
func NewTitle(a api.API, param string) Fragment {
	return &Title{a}
}

// Text implements Fragment.
func (w *Title) Text() (string, string) {
	songlistWidget := w.api.SonglistWidget()
	songlist := songlistWidget.Songlist()
	return songlist.Name(), `title`
}
