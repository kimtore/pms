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

func (w *SonglistWidget) List() *list {
	return w.currentList
}

func (w *SonglistWidget) Songlist() songlist.Songlist {
	return w.List().songlist
}

func (w *SonglistWidget) AutoSetColumnWidths() {
	currentList := w.List()
	for i := range currentList.columns {
		currentList.columns[i].SetWidth(currentList.columns[i].Mid())
	}
	w.expandColumns()
}

func (w *SonglistWidget) SetColumns(cols []string) {
	timer := time.Now()

	currentList := w.List()
	ch := make(chan int, len(currentList.columns))
	currentList.columns = make(columns, len(cols))

	for i := range currentList.columns {
		go func(i int) {
			currentList.columns[i].Tag = cols[i]
			currentList.columns[i].Set(currentList.songlist)
			ch <- 0
		}(i)
	}
	for i := 0; i < len(currentList.columns); i++ {
		<-ch
	}
	w.AutoSetColumnWidths()

	console.Log("Calculated column widths in %s", time.Since(timer).String())

	PostEventListChanged(w)
}

func (w *SonglistWidget) expandColumns() {
	currentList := w.List()

	if len(currentList.columns) == 0 {
		return
	}

	_, _, xmax, _ := w.viewport.GetVisible()
	totalWidth := 0
	poolSize := len(currentList.columns)
	saturation := make([]bool, poolSize)
	for i := range currentList.columns {
		totalWidth += currentList.columns[i].Width()
	}
	for {
		for i := range currentList.columns {
			if totalWidth > xmax {
				return
			}
			if poolSize > 0 && saturation[i] {
				continue
			}
			col := &currentList.columns[i]
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
	currentList := w.List()
	xmax += 1
	style := w.Style("default")
	cursor := false

	for y := ymin; y <= ymax; y++ {

		lineStyled := true
		s := list.Song(y)
		if s == nil {
			// Sometimes happens under race conditions; just abort drawing
			console.Log("Attempting to draw nil song, aborting draw due to possible race condition.")
			return
		}

		// Style based on song's role
		cursor = y == currentList.cursor
		switch {
		case cursor:
			style = w.Style("cursor")
		case w.IndexAtCurrentSong(y):
			style = w.Style("currentSong")
		case currentList.Selected(y):
			style = w.Style("selected")
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
		for col := 0; col < len(currentList.columns); col++ {

			// Convert tag to runes
			key := currentList.columns[col].Tag
			runes := s.Tags[key]
			if !lineStyled {
				style = w.Style(key)
			}

			if col+1 == len(currentList.columns) {
				rightPadding = 0
			}

			strmax := currentList.columns[col].Width()
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
	w.SetCursor(w.List().cursor + i)
}

func (w *SonglistWidget) SetCursor(i int) {
	currentList := w.List()
	currentList.cursor = i
	w.validateCursorList()
	w.expandVisualSelection()
	w.viewport.MakeVisible(0, currentList.cursor)
	w.validateCursorVisible()
	w.PostEventWidgetContent(w)
}

func (w *SonglistWidget) Cursor() int {
	return w.List().cursor
}

func (w *SonglistWidget) CursorSong() *song.Song {
	return w.List().CursorSong()
}

func (w *SonglistWidget) CursorToSong(s *song.Song) error {
	index, err := w.Songlist().Locate(s)
	if err != nil {
		return err
	}
	w.SetCursor(index)
	return nil
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
	currentList := w.List()
	currentList.ymin, currentList.ymax = w.getVisibleBoundaries()
	w.validateCursor(currentList.ymin, currentList.ymax)
}

// validateCursorList makes sure the cursor stays within songlist boundaries.
func (w *SonglistWidget) validateCursorList() {
	ymin, ymax := w.getBoundaries()
	w.validateCursor(ymin, ymax)
}

// validateCursor adjusts the cursor based on minimum and maximum boundaries.
func (w *SonglistWidget) validateCursor(ymin, ymax int) {
	currentList := w.List()
	if currentList.cursor < ymin {
		currentList.cursor = ymin
	}
	if currentList.cursor > ymax {
		currentList.cursor = ymax
	}
}

func (w *SonglistWidget) Resize() {
	w.viewport.Resize(0, 0, -1, -1)
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
	return w.List().columns
}

// Len returns the number of songs in the current songlist.
func (w *SonglistWidget) Len() int {
	return w.Songlist().Len()
}

// SonglistsLen returns the number of indexed songlists.
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

func (w *SonglistWidget) AddSonglist(s songlist.Songlist) {
	list := newList(s)
	w.lists = append(w.lists, list)
	console.Log("Songlist UI: added songlist index %d of type %T at address %p", len(w.lists)-1, s, s)
}

func (w *SonglistWidget) RemoveSonglist(index int) error {
	if err := w.ValidateSonglistIndex(index); err != nil {
		return err
	}
	if index+1 == w.SonglistsLen() {
		w.lists = w.lists[:index]
	} else {
		w.lists = append(w.lists[:index], w.lists[index+1:]...)
	}
	console.Log("Songlist UI: removed songlist index %d", index)
	return nil
}

// ReplaceSonglist replaces an existing songlist with its new version. Checking
// is done on a type-level, so only the queue and library will be replaced.
func (w *SonglistWidget) ReplaceSonglist(s songlist.Songlist) {
	currentList := w.List()

	for i := range w.lists {
		if reflect.TypeOf(w.lists[i].songlist) != reflect.TypeOf(s) {
			continue
		}
		console.Log("Songlist UI: replacing songlist of type %T at %p with new list at %p", s, w.lists[i].songlist, s)
		console.Log("Songlist UI: comparing %p %p", w.lists[i], currentList)

		active := w.lists[i] == currentList
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
	w.ListChanged()
}

func (w *SonglistWidget) ListChanged() {
	//if len(w.currentList.columns) == 0 {
	w.SetColumns(strings.Split(w.options.StringValue("columns"), ",")) // FIXME
	//}
	w.setViewportSize()
	currentList := w.List()
	w.viewport.MakeVisible(0, currentList.ymax)
	w.viewport.MakeVisible(0, currentList.ymin)
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

func (w *SonglistWidget) ValidateSonglistIndex(i int) error {
	if !w.ValidSonglistIndex(i) {
		return fmt.Errorf("Index %d is out of bounds (try between 1 and %d)", i+1, w.SonglistsLen())
	}
	return nil
}

func (w *SonglistWidget) SetSonglistIndex(i int) error {
	console.Log("SetSonglistIndex(%d)", i)
	if err := w.ValidateSonglistIndex(i); err != nil {
		return err
	}
	w.currentListIndex = i
	w.activateList(w.lists[w.currentListIndex])
	w.SetFallbackSonglist(w.List().songlist)
	return nil
}

// Selection returns the current selection as a songlist.Songlist.
func (w *SonglistWidget) Selection() songlist.Songlist {
	list := w.List()
	indices := list.SelectionIndices()
	source := w.Songlist()
	dest := songlist.New()
	for _, i := range indices {
		if song := source.Song(i); song != nil {
			dest.Add(song)
		} else {
			console.Log("SelectionIndices() returned an integer '%d' that resulted in a nil song, ignoring", i)
		}
	}
	return dest
}

// EnableVisualSelection sets start and stop of the visual selection to the
// cursor position.
func (w *SonglistWidget) EnableVisualSelection() {
	list := w.List()
	w.SetVisualSelection(list.cursor, list.cursor, list.cursor)
}

// DisableVisualSelection disables visual selection.
func (w *SonglistWidget) DisableVisualSelection() {
	w.SetVisualSelection(-1, -1, -1)
}

// ClearSelection clears the selection.
func (w *SonglistWidget) ClearSelection() {
	w.List().ClearSelection()
	PostEventModeSync(w, MultibarModeNormal)
}

// SetVisualSelection sets the range of the visual selection. Use negative
// integers to un-select all visually selected songs.
func (w *SonglistWidget) SetVisualSelection(ymin, ymax, ystart int) {
	list := w.List()
	list.SetVisualSelection(ymin, ymax, ystart)
	if list.HasVisualSelection() {
		PostEventModeSync(w, MultibarModeVisual)
	} else {
		PostEventModeSync(w, MultibarModeNormal)
	}
}

// HasVisualSelection returns true if the songlist is in visual selection mode.
func (w *SonglistWidget) HasVisualSelection() bool {
	return w.List().HasVisualSelection()
}

// expandVisualSelection sets the visual selection boundaries from where it
// started to the current cursor position.
func (w *SonglistWidget) expandVisualSelection() {
	list := w.List()
	if !list.HasVisualSelection() {
		return
	}
	ymin, ymax, ystart := list.VisualSelection()
	switch {
	case list.cursor < ystart:
		ymin, ymax = list.cursor, ystart
	case list.cursor > ystart:
		ymin, ymax = ystart, list.cursor
	default:
		ymin, ymax = ystart, ystart
	}
	w.SetVisualSelection(ymin, ymax, ystart)
}
