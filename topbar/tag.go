package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/gdamore/tcell"
)

// Tag draws song tags from the currently playing song.
type Tag struct {
	tag string
	fragment
}

func NewTag(a api.API, param string) Fragment {
	return &Tag{
		param, fragment{api: a},
	}
}

func (w *Tag) Text() ([]rune, tcell.Style) {
	song := w.api.Song()
	if song == nil {
		return []rune(`<none>`), w.Style("tagMissing")
	}
	if text, ok := song.Tags[w.tag]; ok {
		return text, w.Style(w.tag)
	}
	return []rune(`<unknown>`), w.Style("tagMissing")
}

func (w *Tag) Width() int {
	text, _ := w.Text()
	return len(text)
}

func (w *Tag) Draw(x, y int) int {
	text, style := w.Text()
	return w.drawNext(x, y, text, style)
}
