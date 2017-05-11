package widgets

import (
	"fmt"
	"math"
	"reflect"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/utils"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// SonglistWidget is a tcell widget which draws a Songlist on the screen. It
// maintains a list of songlists which can be cycled through.
type SonglistWidget struct {
	api              api.API
	columns          songlist.Columns
	fallbackSonglist songlist.Songlist
	listIndex        int
	songlist         songlist.Songlist
	songlists        []songlist.Songlist

	view     views.View
	viewport views.ViewPort

	style.Styled
	views.WidgetWatchers
}

func NewSonglistWidget(a api.API) (w *SonglistWidget) {
	w = &SonglistWidget{}
	w.songlist = songlist.New()
	w.api = a
	return
}

func (w *SonglistWidget) SetAPI(a api.API) {
	w.api = a
}

func (w *SonglistWidget) Songlist() songlist.Songlist {
	return w.songlist
}

func (w *SonglistWidget) drawNext(x, y, strmin, strmax int, runes []rune, style tcell.Style) int {
	strmin = utils.Min(len(runes), strmin)
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

	//list.Lock()
	//defer list.Unlock()

	_, ymin, xmax, ymax := w.viewport.GetVisible()
	currentSong := w.api.Song()
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
		cursor = y == list.Cursor()
		switch {
		case cursor:
			style = w.Style("cursor")
		case list.IndexAtSong(y, currentSong):
			style = w.Style("currentSong")
		case list.Selected(y):
			style = w.Style("selection")
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
		for col := 0; col < len(w.columns); col++ {

			// Convert tag to runes
			key := w.columns[col].Tag
			runes := s.Tags[key]
			if !lineStyled {
				style = w.Style(key)
			}

			if col+1 == len(w.columns) {
				rightPadding = 0
			}

			strmax := w.columns[col].Width()
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

func (w *SonglistWidget) setViewportSize() {
	x, y := w.Size()
	w.viewport.SetContentSize(x, w.Songlist().Len(), true)
	w.viewport.SetSize(x, utils.Min(y, w.Songlist().Len()))
}

func (w *SonglistWidget) validateViewport() {
	w.validateCursorVisible()
	//w.AutoSetColumnWidths() FIXME
	w.PostEventWidgetContent(w)
}

// validateCursorVisible makes sure the cursor stays within the visible area of the viewport.
func (w *SonglistWidget) validateCursorVisible() {
	ymin, ymax := w.getVisibleBoundaries()
	w.Songlist().ValidateCursor(ymin, ymax)
}

// validateCursorList makes sure the cursor stays within songlist boundaries.
func (w *SonglistWidget) validateCursorList() {
	ymin, ymax := 0, w.Songlist().Len()-1
	w.Songlist().ValidateCursor(ymin, ymax)
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

// SonglistsLen returns the number of indexed songlists.
// FIXME: rename this to Len()
func (w *SonglistWidget) SonglistsLen() int {
	return len(w.songlists)
}

// PositionReadout returns a combination of PositionLongReadout() and PositionShortReadout().
// FIXME: move this into a positionreadout fragment
func (w *SonglistWidget) PositionReadout() string {
	return fmt.Sprintf("%s    %s", w.PositionLongReadout(), w.PositionShortReadout())
}

// PositionLongReadout returns a formatted string containing the visible song
// range as well as the total number of songs.
// FIXME: move this into a positionreadout fragment
func (w *SonglistWidget) PositionLongReadout() string {
	ymin, ymax := w.getVisibleBoundaries()
	return fmt.Sprintf("%d,%d-%d/%d", w.Songlist().Cursor()+1, ymin+1, ymax+1, w.Songlist().Len())
}

// PositionShortReadout returns a percentage indicator on how far the songlist is scrolled.
// FIXME: move this into a positionreadout fragment
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
	w.songlists = append(w.songlists, s)
	console.Log("Songlist UI: added songlist index %d of type %T at address %p", len(w.songlists)-1, s, s)
}

func (w *SonglistWidget) RemoveSonglist(index int) error {
	if err := w.ValidateSonglistIndex(index); err != nil {
		return err
	}
	if index+1 == w.SonglistsLen() {
		w.songlists = w.songlists[:index]
	} else {
		w.songlists = append(w.songlists[:index], w.songlists[index+1:]...)
	}
	console.Log("Songlist UI: removed songlist index %d", index)
	return nil
}

// ReplaceSonglist replaces an existing songlist with its new version. Checking
// is done on a type-level, so only the queue and library will be replaced.
func (w *SonglistWidget) ReplaceSonglist(s songlist.Songlist) {
	for i := range w.songlists {
		if reflect.TypeOf(w.songlists[i]) != reflect.TypeOf(s) {
			continue
		}
		console.Log("Songlist UI: replacing songlist of type %T at %p with new list at %p", s, w.songlists[i], s)
		console.Log("Songlist UI: comparing %p %p", w.songlists[i], w.Songlist())

		active := w.songlists[i] == w.Songlist()
		w.songlists[i] = s

		if active {
			console.Log("Songlist UI: replaced songlist is currently active, switching to new songlist.")
			w.SetSonglist(s)
		}
		return
	}

	console.Log("Songlist UI: adding songlist of type %T at address %p since no similar exists", s, s)
	w.AddSonglist(s)
}

// SetSonglist activates the specified songlist. If the songlist is already in
// the index, that index is activated instead.
func (w *SonglistWidget) SetSonglist(s songlist.Songlist) {
	console.Log("SetSonglist(%T %p)", s, s)
	w.listIndex = -1
	for i, stored := range w.songlists {
		if stored == s {
			w.SetSonglistIndex(i)
			return
		}
	}
	w.activateList(s)
}

// SetFallbackSonglist sets a songlist that should be reverted to in case a search result returns zero results.
func (w *SonglistWidget) SetFallbackSonglist(s songlist.Songlist) {
	console.Log("SetFallbackSonglist(%T %p)", s, s)
	w.fallbackSonglist = s
}

// FallbackSonglist returns the songlist that is reverted to in case of zero results.
func (w *SonglistWidget) FallbackSonglist() songlist.Songlist {
	return w.fallbackSonglist
}

func (w *SonglistWidget) activateList(s songlist.Songlist) {
	//console.Log("activateList(%T %p)", s.songlist, s.songlist)
	w.songlist = s
	w.ListChanged()
}

func (w *SonglistWidget) ListChanged() {
	w.setViewportSize()
	w.viewport.MakeVisible(0, 0)
	w.viewport.MakeVisible(0, w.Songlist().Cursor())
	w.validateViewport()
	PostEventListChanged(w)
}

// SetColumns sets which columns that should be visible
func (w *SonglistWidget) SetColumns(tags []string) {
	w.columns = w.Songlist().Columns(tags)
}

func (w *SonglistWidget) SonglistIndex() (int, error) {
	if !w.ValidSonglistIndex(w.listIndex) {
		return 0, fmt.Errorf("Songlist index is out of range")
	}
	return w.listIndex, nil
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
	w.listIndex = i
	w.activateList(w.songlists[w.listIndex])
	w.SetFallbackSonglist(w.Songlist())
	return nil
}
