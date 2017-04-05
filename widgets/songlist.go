package widgets

import (
	"fmt"
	"math"

	"github.com/ambientsound/pms/songlist"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type SongListWidget struct {
	songlist    *songlist.SongList
	view        views.View
	viewport    views.ViewPort
	cursor      int
	cursorStyle tcell.Style
	columns     []column

	views.WidgetWatchers
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func NewSongListWidget() (w *SongListWidget) {
	w = &SongListWidget{}
	w.songlist = &songlist.SongList{}
	w.columns = make([]column, 0)
	w.cursorStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	return
}

func (w *SongListWidget) SetSongList(s *songlist.SongList) {
	//console.Log("setSongList(%p)", s)
	w.songlist = s
	w.setViewportSize()
	cols := []string{
		"artist",
		"track",
		"title",
		"album",
		"date",
		"time",
	}
	w.SetColumns(cols)
	PostEventListChanged(w)
}

func (w *SongListWidget) AutoSetColumnWidths() {
	for i := range w.columns {
		w.columns[i].SetWidth(w.columns[i].Mid())
	}
	w.expandColumns()
}

func (w *SongListWidget) SetColumns(cols []string) {
	//timer := time.Now()
	ch := make(chan int, len(cols))
	w.columns = make([]column, len(cols))
	for i := range cols {
		go func(i int) {
			w.columns[i].Tag = cols[i]
			w.columns[i].Set(w.songlist)
			ch <- 0
		}(i)
	}
	for i := 0; i < len(cols); i++ {
		<-ch
	}
	w.AutoSetColumnWidths()
	//console.Log("SetColumns on %d songs in %s", w.songlist.Len(), time.Since(timer).String())
}

func (w *SongListWidget) expandColumns() {
	if len(w.columns) == 0 {
		return
	}
	_, _, xmax, _ := w.viewport.GetVisible()
	totalWidth := 0
	poolSize := len(w.columns)
	saturation := make([]bool, poolSize)
	for i := range w.columns {
		totalWidth += w.columns[i].Width()
	}
	for {
		for i := range w.columns {
			if totalWidth > xmax {
				return
			}
			if poolSize > 0 && saturation[i] {
				continue
			}
			col := &w.columns[i]
			if poolSize > 0 && col.Width() > col.MaxWidth() {
				saturation[i] = true
				poolSize--
				continue
			}
			col.SetWidth(col.Width() + 1)
			totalWidth++
		}
	}
}

func (w *SongListWidget) Draw() {
	if w.view == nil || w.songlist == nil {
		return
	}

	ymin, ymax := w.getVisibleBoundaries()
	style := tcell.StyleDefault

	for y := ymin; y <= ymax; y++ {

		s := w.songlist.Songs[y]

		// Style based on song's role
		if y == w.cursor {
			style = w.cursorStyle
		} else {
			style = tcell.StyleDefault
		}
		x := 0
		rightPadding := 1

		// Draw each column separately
		for col := 0; col < len(w.columns); col++ {

			// Convert tag to runes
			str := s.Tags[w.columns[col].Tag]
			runes := []rune(str)

			if col+1 == len(w.columns) {
				rightPadding = 0
			}
			strmin := min(len(runes), w.columns[col].Width()-rightPadding)
			strmax := w.columns[col].Width()
			n := 0
			for n < strmin {
				w.viewport.SetContent(x, y, runes[n], nil, style)
				n++
				x++
			}
			for n < strmax {
				w.viewport.SetContent(x, y, ' ', nil, style)
				n++
				x++
			}
		}
	}
	w.PostEventWidgetContent(w)
	PostEventScroll(w)
}

func (w *SongListWidget) getVisibleBoundaries() (ymin, ymax int) {
	_, ymin, _, ymax = w.viewport.GetVisible()
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
	x, y := w.Size()
	w.viewport.Resize(0, 0, -1, -1)
	w.viewport.SetContentSize(x, w.songlist.Len(), true)
	w.viewport.SetSize(x, min(y, w.songlist.Len()))
	w.validateCursorVisible()
	w.AutoSetColumnWidths()
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
	w.setViewportSize()
	w.PostEventWidgetResize(w)
}

func (w *SongListWidget) HandleEvent(ev tcell.Event) bool {
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
	w.view = v
	w.viewport.SetView(w.view)
}

func (w *SongListWidget) Size() (int, int) {
	return w.view.Size()
}

func (w *SongListWidget) Name() string {
	return w.songlist.Name
}

func (w *SongListWidget) Columns() []column {
	return w.columns
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
