package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/style"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// fragment is a useful base class for implementers of Fragment.
type fragment struct {
	api  api.API
	view views.View
	style.Styled
}

func (w *fragment) drawNext(x, y int, runes []rune, style tcell.Style) int {
	strlen := 0
	for p, r := range runes {
		w.view.SetContent(x+p, y, r, nil, style)
		strlen++
	}
	return x + strlen
}

func (w *fragment) drawNextString(x, y int, s string, style tcell.Style) int {
	return w.drawNext(x, y, []rune(s), style)
}

func (w *fragment) drawNextChar(x, y int, r rune, style tcell.Style) int {
	w.view.SetContent(x, y, r, nil, style)
	return x + 1
}

func (w *fragment) SetDirty(dirty bool) {
}

func (w *fragment) SetView(v views.View) {
	w.view = v
}
