package topbar

import (
	"github.com/gdamore/tcell"
)

// tag implements common functions for drawing song tags.
type tag struct {
	tag string
	fragment
}

func (w *tag) Text() ([]rune, tcell.Style) {
	song := w.api.Song()
	if song == nil {
		return []rune(`<none>`), w.Style("tagMissing")
	}
	if text, ok := song.Tags[w.tag]; ok {
		return text, w.Style(w.tag)
	}
	return []rune(`<unknown>`), w.Style("tagMissing")
}

func (w *tag) Width() int {
	text, _ := w.Text()
	return len(text)
}

func (w *tag) Draw(x, y int) int {
	text, style := w.Text()
	return w.drawNext(x, y, text, style)
}
