package widgets

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type list struct {
	songlist songlist.Songlist
	cursor   int
	columns  columns
	ymin     int
	ymax     int
}

type SonglistWidget struct {
	currentListIndex int
	currentList      *list
	lists            []*list
	fallbackSonglist songlist.Songlist
	currentSong      song.Song
	view             views.View
	viewport         views.ViewPort
	options          *options.Options

	widget
	views.WidgetWatchers
}

func newList(s songlist.Songlist) *list {
	l := &list{}
	l.songlist = s
	l.columns = make(columns, 0)
	return l
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

func NewSonglistWidget(o *options.Options) (w *SonglistWidget) {
	w = &SonglistWidget{}
	w.options = o
	w.currentList = newList(songlist.New())
	return
}

func (w *SonglistWidget) Songlist() songlist.Songlist {
	return w.currentList.songlist
}

func (w *SonglistWidget) AutoSetColumnWidths() {
	for i := range w.currentList.columns {
		w.currentList.columns[i].SetWidth(w.currentList.columns[i].Mid())
	}
	w.expandColumns()
}

func (w *SonglistWidget) SetColumns(cols []string) {
	timer := time.Now()

	ch := make(chan int, len(w.currentList.columns))
	w.currentList.columns = make(columns, len(cols))

	for i := range w.currentList.columns {
		go func(i int) {
			w.currentList.columns[i].Tag = cols[i]
			w.currentList.columns[i].Set(w.currentList.songlist)
			ch <- 0
		}(i)
	}
	for i := 0; i < len(w.currentList.columns); i++ {
		<-ch
	}
	w.AutoSetColumnWidths()

	console.Log("Calculated column widths in %s", time.Since(timer).String())

	PostEventListChanged(w)
}

func (w *SonglistWidget) expandColumns() {
	if len(w.currentList.columns) == 0 {
		return
	}
	_, _, xmax, _ := w.viewport.GetVisible()
	totalWidth := 0
	poolSize := len(w.currentList.columns)
	saturation := make([]bool, poolSize)
	for i := range w.currentList.columns {
		totalWidth += w.currentList.columns[i].Width()
	}
	for {
		for i := range w.currentList.columns {
			if totalWidth > xmax {
				return
			}
			if poolSize > 0 && saturation[i] {
				continue
			}
			col := &w.currentList.columns[i]
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

func (w *SonglistWidget) drawNext(x, y, strmin, strmax int, runes []rune, style tcell.Style) int {
	strmin = min(len(runes), strmin)
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
	return x
}

func (w *SonglistWidget) drawOneTagLine(x, y, xmax int, s *song.Song, tag string, defaultStyle string, style tcell.Style, lineStyled bool) int {
	if !lineStyled {
		style = w.Style(defaultStyle)
	}

	runes := s.Tags[tag]
	strmin := len(runes)

	return w.drawNext(x, y, strmin, xmax+1, runes, style)
}

func (w *SonglistWidget) Draw() {
	list := w.Songlist()
	if w.view == nil || list == nil || list.Songs() == nil {
		return
	}

	list.Lock()
	defer list.Unlock()

	_, ymin, xmax, ymax := w.viewport.GetVisible()
	xmax += 1
	style := w.Style("default")
	lineStyled := false
	cursor := false

	for y := ymin; y <= ymax; y++ {

		s := list.Song(y)

		// Style based on song's role
		cursor = y == w.currentList.cursor
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

		// If all essential tags are missing, draw only the filename
		if !s.HasOneOfTags("artist", "album", "title") {
			w.drawOneTagLine(x, y, xmax+1, s, `file`, `allTagsMissing`, style, lineStyled)
			continue
		}

		// If most essential tags are missing, but the title is present, draw only the title.
		if !s.HasOneOfTags("artist", "album") {
			w.drawOneTagLine(x, y, xmax+1, s, `title`, `mostTagsMissing`, style, lineStyled)
			continue
		}

		// Draw each column separately
		for col := 0; col < len(w.currentList.columns); col++ {

			// Convert tag to runes
			key := w.currentList.columns[col].Tag
			runes := s.Tags[key]
			if !lineStyled {
				style = w.Style(key)
			}

			if col+1 == len(w.currentList.columns) {
				rightPadding = 0
			}

			strmax := w.currentList.columns[col].Width()
			strmin := strmax - rightPadding

			x = w.drawNext(x, y, strmin, strmax, runes, style)
		}
	}
	w.PostEventWidgetContent(w)
	PostEventScroll(w)
}

func (w *SonglistWidget) getVisibleBoundaries() (ymin, ymax int) {
	_, ymin, _, ymax = w.viewport.GetVisible()
	return
}

func (w *SonglistWidget) Width() int {
	_, _, xmax, _ := w.viewport.GetVisible()
	return xmax
}

func (w *SonglistWidget) getBoundaries() (ymin, ymax int) {
	return 0, w.Songlist().Len() - 1
}

func (w *SonglistWidget) setViewportSize() {
	x, y := w.Size()
	w.viewport.Resize(0, 0, -1, -1)
	w.viewport.SetContentSize(x, w.Songlist().Len(), true)
	w.viewport.SetSize(x, min(y, w.Songlist().Len()))
}

func (w *SonglistWidget) validateViewport() {
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
	w.SetCursor(w.currentList.cursor + i)
}

func (w *SonglistWidget) SetCursor(i int) {
	w.currentList.cursor = i
	w.validateCursorList()
	w.viewport.MakeVisible(0, w.currentList.cursor)
	w.validateCursorVisible()
	w.PostEventWidgetContent(w)
}

func (w *SonglistWidget) Cursor() int {
	return w.currentList.cursor
}

func (w *SonglistWidget) CursorSong() *song.Song {
	return w.Songlist().Song(w.currentList.cursor)
}

func (w *SonglistWidget) SetCurrentSong(s *song.Song) {
	if s != nil {
		w.currentSong = *s
	} else {
		w.currentSong = song.Song{}
	}
}

func (w *SonglistWidget) IndexAtCurrentSong(i int) bool {
	s := w.Songlist().Song(i)
	if s == nil {
		return false
	}
	if songlist.IsQueue(w.Songlist()) {
		return s.ID == w.currentSong.ID
	} else {
		return s.StringTags["file"] == w.currentSong.StringTags["file"]
	}
}

// validateCursorVisible makes sure the cursor stays within the visible area of the viewport.
func (w *SonglistWidget) validateCursorVisible() {
	w.currentList.ymin, w.currentList.ymax = w.getVisibleBoundaries()
	w.validateCursor(w.currentList.ymin, w.currentList.ymax)
}

// validateCursorList makes sure the cursor stays within songlist boundaries.
func (w *SonglistWidget) validateCursorList() {
	ymin, ymax := w.getBoundaries()
	w.validateCursor(ymin, ymax)
}

// validateCursor adjusts the cursor based on minimum and maximum boundaries.
func (w *SonglistWidget) validateCursor(ymin, ymax int) {
	if w.currentList.cursor < ymin {
		w.currentList.cursor = ymin
	}
	if w.currentList.cursor > ymax {
		w.currentList.cursor = ymax
	}
}

func (w *SonglistWidget) Resize() {
	w.setViewportSize()
	w.validateViewport()
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
	return w.Songlist().Name()
}

func (w *SonglistWidget) Columns() []column {
	return w.currentList.columns
}

func (w *SonglistWidget) Len() int {
	return w.Songlist().Len()
}

func (w *SonglistWidget) SonglistsLen() int {
	return len(w.lists)
}

// PositionReadout returns a combination of PositionLongReadout() and PositionShortReadout().
func (w *SonglistWidget) PositionReadout() string {
	return fmt.Sprintf("%s    %s", w.PositionLongReadout(), w.PositionShortReadout())
}

// PositionLongReadout returns a formatted string containing the visible song
// range as well as the total number of songs.
func (w *SonglistWidget) PositionLongReadout() string {
	ymin, ymax := w.getVisibleBoundaries()
	return fmt.Sprintf("%d,%d-%d/%d", w.Cursor()+1, ymin+1, ymax+1, w.Songlist().Len())
}

// PositionShortReadout returns a percentage indicator on how far the songlist is scrolled.
func (w *SonglistWidget) PositionShortReadout() string {
	ymin, ymax := w.getVisibleBoundaries()
	if ymin == 0 && ymax+1 == w.Songlist().Len() {
		return `All`
	}
	if ymin == 0 {
		return `Top`
	}
	if ymax+1 == w.Songlist().Len() {
		return `Bot`
	}
	fraction := float64(float64(ymin) / float64(w.Songlist().Len()))
	percent := int(math.Floor(fraction * 100))
	return fmt.Sprintf("%2d%%", percent)
}

//

func (w *SonglistWidget) AddSonglist(s songlist.Songlist) {
	list := newList(s)
	w.lists = append(w.lists, list)
	console.Log("Songlist UI: added songlist index %d of type %T at address %p", len(w.lists)-1, s, s)
}

// ReplaceSonglist replaces an existing songlist with its new version. Checking
// is done on a type-level, so only the queue and library will be replaced.
func (w *SonglistWidget) ReplaceSonglist(s songlist.Songlist) {
	for i := range w.lists {
		if reflect.TypeOf(w.lists[i].songlist) != reflect.TypeOf(s) {
			continue
		}
		console.Log("Songlist UI: replacing songlist of type %T at %p with new list at %p", s, w.lists[i].songlist, s)
		console.Log("Songlist UI: comparing %p %p", w.lists[i], w.currentList)

		active := w.lists[i] == w.currentList
		w.lists[i].songlist = s

		if active {
			console.Log("Songlist UI: replaced songlist is currently active, switching to new songlist.")
			w.SetSonglist(s)
		}
		return
	}

	console.Log("Songlist UI: adding songlist of type %T at address %p since no similar exists", s, s)
	w.AddSonglist(s)
}

func (w *SonglistWidget) SetSonglist(s songlist.Songlist) {
	console.Log("SetSonglist(%T %p)", s, s)
	w.currentListIndex = -1
	for i, stored := range w.lists {
		if stored.songlist == s {
			w.SetSonglistIndex(i)
			return
		}
	}
	w.activateList(newList(s))
}

// SetFallbackSonglist sets a songlist that should be reverted to in case a search result returns zero results.
func (w *SonglistWidget) SetFallbackSonglist(s songlist.Songlist) {
	console.Log("SetFallbackSonglist(%T %p)", s, s)
	w.fallbackSonglist = s
}

func (w *SonglistWidget) FallbackSonglist() songlist.Songlist {
	return w.fallbackSonglist
}

func (w *SonglistWidget) activateList(s *list) {
	console.Log("activateList(%T %p)", s.songlist, s.songlist)
	w.currentList = s
	//if len(w.currentList.columns) == 0 {
	w.SetColumns(strings.Split(w.options.StringValue("columns"), ",")) // FIXME
	//}
	w.setViewportSize()
	//console.Log("Calling MakeVisible(%d), MakeVisible(%d)", w.currentList.ymax, w.currentList.ymin)
	w.viewport.MakeVisible(0, w.currentList.ymax)
	w.viewport.MakeVisible(0, w.currentList.ymin)
	w.validateViewport()
	PostEventListChanged(w)
}

func (w *SonglistWidget) SonglistIndex() (int, error) {
	if !w.ValidSonglistIndex(w.currentListIndex) {
		return 0, fmt.Errorf("Songlist index is out of range")
	}
	return w.currentListIndex, nil
}

func (w *SonglistWidget) ValidSonglistIndex(i int) bool {
	return i >= 0 && i < w.SonglistsLen()
}

func (w *SonglistWidget) SetSonglistIndex(i int) error {
	console.Log("SetSonglistIndex(%d)", i)
	if !w.ValidSonglistIndex(i) {
		return fmt.Errorf("Index %d is out of bounds (try between 1 and %d)", i+1, w.SonglistsLen())
	}
	w.currentListIndex = i
	w.activateList(w.lists[w.currentListIndex])
	w.SetFallbackSonglist(w.currentList.songlist)
	return nil
}
