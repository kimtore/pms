package widgets

import (
	"github.com/ambientsound/pms/songlist"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"

	"fmt"
	"math"
)

type SongListWidget struct {
	songlist    *songlist.SongList
	view        views.View
	viewport    views.ViewPort
	cursor      int
	cursorStyle tcell.Style

	views.WidgetWatchers
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func NewSongListWidget() (w *SongListWidget) {
	w = &SongListWidget{}
	w.songlist = &songlist.SongList{}
	w.cursorStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	return
}

func (w *SongListWidget) SetSongList(s *songlist.SongList) {
	//console.Log("setSongList(%p)", s)
	w.songlist = s
	w.setViewportSize()
	PostEventListChanged(w)
}

func (w *SongListWidget) Draw() {
	//console.Log("Draw() with view=%p", w.view)
	if w.view == nil || w.songlist == nil {
		return
	}
	ymin, ymax := w.getVisibleBoundaries()
	style := tcell.StyleDefault
	for y := ymin; y <= ymax; y++ {
		s := w.songlist.Songs[y]
		str := s.Tags["file"]
		if y == w.cursor {
			style = w.cursorStyle
		} else {
			style = tcell.StyleDefault
		}
		for x := 0; x < len(str); x++ {
			w.viewport.SetContent(x, y, rune(str[x]), nil, style)
		}
	}
	w.PostEventWidgetContent(w)
	PostEventScroll(w)
}

func (w *SongListWidget) getVisibleBoundaries() (ymin, ymax int) {
	_, ymin, _, ymax = w.viewport.GetVisible()
	//console.Log("GetVisible() says ymin %d, ymax %d", ymin, ymax)
	return
}

func (w *SongListWidget) getBoundaries() (ymin, ymax int) {
	return 0, w.songlist.Len() - 1
}

func (w *SongListWidget) validateCursorVisible() {
	ymin, ymax := w.getVisibleBoundaries()
	w.validateCursor(ymin, ymax)
}

func (w *SongListWidget) validateCursorList() {
	ymin, ymax := w.getBoundaries()
	w.validateCursor(ymin, ymax)
}

func (w *SongListWidget) validateCursor(ymin, ymax int) {
	if w.cursor < ymin {
		w.cursor = ymin
	}
	if w.cursor > ymax {
		w.cursor = ymax
	}
}

func (w *SongListWidget) setViewportSize() {
	//console.Log("setViewportSize()")
	x, y := w.Size()
	w.viewport.Resize(0, 0, -1, -1)
	w.viewport.SetContentSize(x, w.songlist.Len(), true)
	//console.Log("SetContentSize(%d, %d)", x, w.songlist.Len())
	w.viewport.SetSize(x, min(y, w.songlist.Len()))
	//console.Log("SetSize(%d, %d)", x, min(y, w.songlist.Len()))
	w.validateCursorVisible()
	w.PostEventWidgetContent(w)
}

func (w *SongListWidget) MoveCursorUp(i int) {
	w.MoveCursor(-i)
}

func (w *SongListWidget) MoveCursorDown(i int) {
	w.MoveCursor(i)
}

func (w *SongListWidget) MoveCursor(i int) {
	w.SetCursor(w.cursor + i)
}

func (w *SongListWidget) SetCursor(i int) {
	w.cursor = i
	w.validateCursorList()
	w.viewport.MakeVisible(0, i)
	w.validateCursorVisible()
	w.PostEventWidgetContent(w)
}

func (w *SongListWidget) Cursor() int {
	return w.cursor
}

func (w *SongListWidget) Resize() {
	//console.Log("Resize()")
	w.setViewportSize()
}

func (w *SongListWidget) HandleEvent(ev tcell.Event) bool {
	//console.Log("HandleEvent()")
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyUp:
			w.MoveCursorUp(1)
			return true
		case tcell.KeyDown:
			w.MoveCursorDown(1)
			return true
		case tcell.KeyPgUp:
			_, y := w.Size()
			w.MoveCursorUp(y)
			return true
		case tcell.KeyPgDn:
			_, y := w.Size()
			w.MoveCursorDown(y)
			return true
		case tcell.KeyHome:
			w.SetCursor(0)
			return true
		case tcell.KeyEnd:
			w.SetCursor(w.songlist.Len() - 1)
			return true
		}
	}
	return false
}

func (w *SongListWidget) SetView(v views.View) {
	//console.Log("SetView(%p)", v)
	w.view = v
	w.viewport.SetView(w.view)
}

func (w *SongListWidget) Size() (int, int) {
	//x, y := w.view.Size()
	//console.Log("Size() returns %d, %d", x, y)
	return w.view.Size()
}

// PositionReadout returns a combination of PositionLongReadout() and PositionShortReadout().
func (w *SongListWidget) PositionReadout() string {
	return fmt.Sprintf("%s    %s", w.PositionLongReadout(), w.PositionShortReadout())
}

// PositionLongReadout returns a formatted string containing the visible song
// range as well as the total number of songs.
func (w *SongListWidget) PositionLongReadout() string {
	ymin, ymax := w.getVisibleBoundaries()
	return fmt.Sprintf("%d,%d-%d/%d", w.Cursor()+1, ymin+1, ymax+1, w.songlist.Len())
}

// PositionShortReadout returns a percentage indicator on how far the songlist is scrolled.
func (w *SongListWidget) PositionShortReadout() string {
	ymin, ymax := w.getVisibleBoundaries()
	if ymin == 0 && ymax+1 == w.songlist.Len() {
		return `All`
	}
	if ymin == 0 {
		return `Top`
	}
	if ymax+1 == w.songlist.Len() {
		return `Bot`
	}
	fraction := float64(float64(ymin) / float64(w.songlist.Len()))
	percent := int(math.Floor(fraction * 100))
	return fmt.Sprintf("%2d%%", percent)
}
