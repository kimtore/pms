package widgets

import (
	"fmt"
	"math"
	"time"

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
	api     api.API
	columns songlist.Columns

	view     views.View
	viewport views.ViewPort
	lastDraw time.Time

	style.Styled
	views.WidgetWatchers
}

func NewSonglistWidget(a api.API) (w *SonglistWidget) {
	return &SonglistWidget{
		api: a,
	}
}

func (w *SonglistWidget) SetAPI(a api.API) {
	w.api = a
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

func (w *SonglistWidget) Panel() *songlist.Collection {
	return w.api.Db().Panel()
}

func (w *SonglistWidget) List() songlist.Songlist {
	return w.Panel().Current()
}

func (w *SonglistWidget) Draw() {
	console.Log("Draw() in songlist widget")
	list := w.List()
	if w.view == nil || list == nil || list.Songs() == nil {
		console.Log(".. BUG: nil list, aborting draw!")
		return
	}

	// Check if the current panel's songlist has changed.
	if w.Panel().Updated().After(w.lastDraw) {
		w.setViewportSize()
		PostEventListChanged(w)
	}

	w.lastDraw = time.Now()

	w.validateViewport()

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
			key := w.columns[col].Tag()
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

func (w *SonglistWidget) GetVisibleBoundaries() (ymin, ymax int) {
	_, ymin, _, ymax = w.viewport.GetVisible()
	return
}

// Width returns the widget width.
func (w *SonglistWidget) Width() int {
	_, _, xmax, _ := w.viewport.GetVisible()
	return xmax
}

// Height returns the widget height.
func (w *SonglistWidget) Height() int {
	_, ymin, _, ymax := w.viewport.GetVisible()
	return ymax - ymin
}

func (w *SonglistWidget) setViewportSize() {
	x, y := w.Size()
	w.viewport.SetContentSize(x, w.List().Len(), true)
	w.viewport.SetSize(x, utils.Min(y, w.List().Len()))
	w.validateViewport()
}

// validateViewport moves the visible viewport so that the cursor is made visible.
// If the 'center' option is enabled, the viewport is centered on the cursor.
func (w *SonglistWidget) validateViewport() {
	list := w.List()
	cursor := list.Cursor()

	// Make the cursor visible
	if !w.api.Options().BoolValue("center") {
		w.viewport.MakeVisible(0, cursor)
		return
	}

	// If 'center' is on, make the cursor centered.
	half := w.Height() / 2
	min := utils.Max(0, cursor-half)
	max := utils.Min(list.Len()-1, cursor+half)
	w.viewport.MakeVisible(0, min)
	w.viewport.MakeVisible(0, max)
}

func (w *SonglistWidget) Resize() {
	w.viewport.Resize(0, 0, -1, -1)
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
	return w.List().Name()
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
	ymin, ymax := w.GetVisibleBoundaries()
	return fmt.Sprintf("%d,%d-%d/%d", w.List().Cursor()+1, ymin+1, ymax+1, w.List().Len())
}

// PositionShortReadout returns a percentage indicator on how far the songlist is scrolled.
// FIXME: move this into a positionreadout fragment
func (w *SonglistWidget) PositionShortReadout() string {
	ymin, ymax := w.GetVisibleBoundaries()
	if ymin == 0 && ymax+1 == w.List().Len() {
		return `All`
	}
	if ymin == 0 {
		return `Top`
	}
	if ymax+1 == w.List().Len() {
		return `Bot`
	}
	fraction := float64(float64(ymin) / float64(w.List().Len()))
	percent := int(math.Floor(fraction * 100))
	return fmt.Sprintf("%2d%%", percent)
}

// SetColumns sets which columns that should be visible
func (w *SonglistWidget) SetColumns(tags []string) {
	xmax, _ := w.Size()
	w.columns = w.List().Columns(tags)
	w.columns.Expand(xmax)
	//console.Log("SetColumns(%v) yields %+v", tags, w.columns)
}

// ScrollViewport scrolls the viewport by delta rows, as far as possible.
// If movecursor is false, the cursor is kept pointing at the same song where
// possible. If true, the cursor is moved delta rows.
func (w *SonglistWidget) ScrollViewport(delta int, movecursor bool) {
	// Do nothing if delta is zero
	if delta == 0 {
		return
	}

	if delta < 0 {
		w.viewport.ScrollUp(-delta)
	} else {
		w.viewport.ScrollDown(delta)
	}

	if movecursor {
		w.List().MoveCursor(delta)
	}

	w.validateCursor()
}

// validateCursor ensures the cursor is within the allowable area without moving
// the viewport.
func (w *SonglistWidget) validateCursor() {
	ymin, ymax := w.GetVisibleBoundaries()
	list := w.List()
	cursor := list.Cursor()

	if w.api.Options().BoolValue("center") {
		// When 'center' is on, move cursor to the centre of the viewport
		target := cursor
		lowerbound := (ymin + ymax) / 2
		upperbound := lowerbound
		if ymin <= 0 {
			// We are scrolled to the top, so the cursor is allowed to go above
			// the middle of the viewport
			lowerbound = 0
		}
		if ymax >= list.Len()-1 {
			// We are scrolled to the bottom, so the cursor is allowed to go
			// below the middle of the viewport
			upperbound = list.Len() - 1
		}
		if target < lowerbound {
			target = lowerbound
		}
		if target > upperbound {
			target = upperbound
		}
		list.SetCursor(target)
	} else {
		// When 'center' is off, move cursor into the viewport
		if cursor < ymin {
			list.SetCursor(ymin)
		} else if cursor > ymax {
			list.SetCursor(ymax)
		}
	}
}
