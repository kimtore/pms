package widgets

import (
	"fmt"
	"math"

	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type SonglistWidget struct {
	songlist    songlist.Songlist
	currentSong song.Song
	view        views.View
	viewport    views.ViewPort
	cursor      int
	columns     []column

	widget
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

func NewSonglistWidget() (w *SonglistWidget) {
	w = &SonglistWidget{}
	w.songlist = &songlist.BaseSonglist{}
	w.columns = make([]column, 0)
	return
}

func (w *SonglistWidget) SetSonglist(s songlist.Songlist) {
	//console.Log("setSonglist(%p)", s)
	w.songlist = s
	w.setViewportSize()
	PostEventListChanged(w)
}

func (w *SonglistWidget) AutoSetColumnWidths() {
	for i := range w.columns {
		w.columns[i].SetWidth(w.columns[i].Mid())
	}
	w.expandColumns()
}

func (w *SonglistWidget) SetColumns(cols []string) {
	ch := make(chan int, len(w.columns))
	w.columns = make([]column, len(cols))

	for i := range w.columns {
		go func(i int) {
			w.columns[i].Tag = cols[i]
			w.columns[i].Set(w.songlist)
			ch <- 0
		}(i)
	}
	for i := 0; i < len(w.columns); i++ {
		<-ch
	}
	w.AutoSetColumnWidths()
	PostEventListChanged(w)
}

func (w *SonglistWidget) expandColumns() {
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

func (w *SonglistWidget) Draw() {
	if w.view == nil || w.songlist == nil || w.songlist.Songs == nil {
		return
	}

	ymin, ymax := w.getVisibleBoundaries()
	style := w.Style("default")
	lineStyled := false
	cursor := false

	for y := ymin; y <= ymax; y++ {

		s := w.songlist.Song(y)

		// Style based on song's role
		cursor = y == w.cursor
		switch {
		case cursor:
			style = w.Style("cursor")
			lineStyled = true
		case w.IndexAtCurrentSong(y):
			style = w.Style("currentSong")
			lineStyled = true
		default:
			style = w.Style("default")
			lineStyled = false
		}

		x := 0
		rightPadding := 1

		// Draw each column separately
		for col := 0; col < len(w.columns); col++ {

			// Convert tag to runes
			key := w.columns[col].Tag
			str := s.Tags[key]
			if !lineStyled {
				style = w.Style(key)
			}
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

func (w *SonglistWidget) getVisibleBoundaries() (ymin, ymax int) {
	_, ymin, _, ymax = w.viewport.GetVisible()
	return
}

func (w *SonglistWidget) getBoundaries() (ymin, ymax int) {
	return 0, w.songlist.Len() - 1
}

func (w *SonglistWidget) setViewportSize() {
	x, y := w.Size()
	w.viewport.Resize(0, 0, -1, -1)
	w.viewport.SetContentSize(x, w.songlist.Len(), true)
	w.viewport.SetSize(x, min(y, w.songlist.Len()))
	w.validateCursorVisible()
	w.AutoSetColumnWidths()
	w.PostEventWidgetContent(w)
}

func (w *SonglistWidget) MoveCursorUp(i int) {
	w.MoveCursor(-i)
}

func (w *SonglistWidget) MoveCursorDown(i int) {
	w.MoveCursor(i)
}

func (w *SonglistWidget) MoveCursor(i int) {
	w.SetCursor(w.cursor + i)
}

func (w *SonglistWidget) SetCursor(i int) {
	w.cursor = i
	w.validateCursorList()
	w.viewport.MakeVisible(0, w.cursor)
	w.validateCursorVisible()
	w.PostEventWidgetContent(w)
}

func (w *SonglistWidget) Cursor() int {
	return w.cursor
}

func (w *SonglistWidget) CursorSong() *song.Song {
	return w.songlist.Song(w.cursor)
}

func (w *SonglistWidget) SetCurrentSong(s *song.Song) {
	if s != nil {
		w.currentSong = *s
	} else {
		w.currentSong = song.Song{}
	}
}

func (w *SonglistWidget) IndexAtCurrentSong(i int) bool {
	if s := w.songlist.Song(i); s != nil {
		return s.TagString("file") == w.currentSong.TagString("file")
	}
	return false
}

// validateCursorVisible makes sure the cursor stays within the visible area of the viewport.
func (w *SonglistWidget) validateCursorVisible() {
	ymin, ymax := w.getVisibleBoundaries()
	w.validateCursor(ymin, ymax)
}

// validateCursorList makes sure the cursor stays within songlist boundaries.
func (w *SonglistWidget) validateCursorList() {
	ymin, ymax := w.getBoundaries()
	w.validateCursor(ymin, ymax)
}

// validateCursor adjusts the cursor based on minimum and maximum boundaries.
func (w *SonglistWidget) validateCursor(ymin, ymax int) {
	if w.cursor < ymin {
		w.cursor = ymin
	}
	if w.cursor > ymax {
		w.cursor = ymax
	}
}

func (w *SonglistWidget) Resize() {
	w.setViewportSize()
	w.PostEventWidgetResize(w)
}

func (m *SonglistWidget) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *SonglistWidget) SetView(v views.View) {
	w.view = v
	w.viewport.SetView(w.view)
}

func (w *SonglistWidget) Size() (int, int) {
	return w.view.Size()
}

func (w *SonglistWidget) Name() string {
	return w.songlist.Name()
}

func (w *SonglistWidget) Columns() []column {
	return w.columns
}

func (w *SonglistWidget) Len() int {
	return w.songlist.Len()
}

// PositionReadout returns a combination of PositionLongReadout() and PositionShortReadout().
func (w *SonglistWidget) PositionReadout() string {
	return fmt.Sprintf("%s    %s", w.PositionLongReadout(), w.PositionShortReadout())
}

// PositionLongReadout returns a formatted string containing the visible song
// range as well as the total number of songs.
func (w *SonglistWidget) PositionLongReadout() string {
	ymin, ymax := w.getVisibleBoundaries()
	return fmt.Sprintf("%d,%d-%d/%d", w.Cursor()+1, ymin+1, ymax+1, w.songlist.Len())
}

// PositionShortReadout returns a percentage indicator on how far the songlist is scrolled.
func (w *SonglistWidget) PositionShortReadout() string {
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
