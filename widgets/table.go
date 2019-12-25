package widgets

import (
	"fmt"
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/log"
	"math"
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/utils"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type column struct {
	col          list.Column
	key          string
	rightPadding int
	width        int
}

// Table is a tcell widget which draws a gridded table.
// maintains a list of songlists which can be cycled through.
type Table struct {
	api            api.API
	visibleColumns []string
	columns        []column
	list           list.List

	view     views.View
	viewport views.ViewPort
	lastDraw time.Time

	style.Styled
	views.WidgetWatchers
}

func NewTable(a api.API) *Table {
	return &Table{
		api: a,
	}
}

func (w *Table) drawNext(x, y, strmin, strmax int, runes []rune, style tcell.Style) int {
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

func (w *Table) drawOneTagLine(x, y, xmax int, s *song.Song, tag string, defaultStyle string, style tcell.Style, lineStyled bool) int {
	if !lineStyled {
		style = w.Style(defaultStyle)
	}

	runes := s.Tags[tag]
	strmin := len(runes)

	return w.drawNext(x, y, strmin, xmax+1, runes, style)
}

func (w *Table) SetList(lst list.List) {
	w.list = lst
	w.SetColumns(lst.ColumnNames())
}

func (w *Table) Draw() {
	log.Debugf("table widget: draw")

	w.SetStylesheet(w.api.Styles())

	// Make sure that the viewport matches the list size.
	w.setViewportSize()
	log.Debugf("%#v", w.viewport)

	// Update draw time
	w.lastDraw = time.Now()

	_, ymin, xmax, ymax := w.viewport.GetVisible()
	xmax += 1
	st := w.Style("default")
	cursor := false

	for y := ymin; y <= ymax; y++ {

		row := w.list.Row(y)
		if row == nil {
			panic("nil row")
		}

		lineStyled := true

		// Style based on song's role
		cursor = y == w.list.Cursor()
		switch {
		case cursor:
			st = w.Style("cursor")
			/*
			// FIXME: dealing with current song?
		case w.list.IndexAtSong(y, currentSong):
			st = w.Style("currentSong")
			*/
		case w.list.Selected(y):
			st = w.Style("selection")
		default:
			st = w.Style("default")
			lineStyled = false
		}

		x := 0

		// Draw each column separately
		for col, name := range w.visibleColumns {

			runes := []rune(row[name])
			if !lineStyled {
				st = w.Style(name)
			}

			strmax := w.columns[col].width
			strmin := strmax - w.columns[col].rightPadding

			x = w.drawNext(x, y, strmin, strmax, runes, st)
		}
	}
}

func (w *Table) GetVisibleBoundaries() (ymin, ymax int) {
	_, ymin, _, ymax = w.viewport.GetVisible()
	return
}

// Width returns the widget width.
func (w *Table) Width() int {
	_, _, xmax, _ := w.viewport.GetVisible()
	return xmax
}

// Height returns the widget height.
func (w *Table) Height() int {
	_, ymin, _, ymax := w.viewport.GetVisible()
	return ymax - ymin
}

func (w *Table) setViewportSize() {
	x, y := w.Size()
	w.viewport.SetContentSize(x, w.list.Len(), true)
	w.viewport.SetSize(x, utils.Min(y, w.list.Len()))
	w.validateViewport()
}

// validateViewport moves the visible viewport so that the cursor is made visible.
// If the 'center' option is enabled, the viewport is centered on the cursor.
func (w *Table) validateViewport() {
	cursor := w.list.Cursor()

	// Make the cursor visible
	if !w.api.Options().BoolValue("center") {
		w.viewport.MakeVisible(0, cursor)
		return
	}

	// If 'center' is on, make the cursor centered.
	half := w.Height() / 2
	min := utils.Max(0, cursor-half)
	max := utils.Min(w.list.Len()-1, cursor+half)
	w.viewport.MakeVisible(0, min)
	w.viewport.MakeVisible(0, max)
}

func (w *Table) Resize() {
	w.SetColumns(w.visibleColumns)
}

func (w *Table) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *Table) SetView(v views.View) {
	w.view = v
	w.viewport.SetView(w.view)
}

func (w *Table) Size() (int, int) {
	return w.view.Size()
}

func (w *Table) Name() string {
	return w.list.Name()
}

// PositionReadout returns a combination of PositionLongReadout() and PositionShortReadout().
// FIXME: move this into a positionreadout fragment
func (w *Table) PositionReadout() string {
	return fmt.Sprintf("%s    %s", w.PositionLongReadout(), w.PositionShortReadout())
}

// PositionLongReadout returns a formatted string containing the visible song
// range as well as the total number of songs.
// FIXME: move this into a positionreadout fragment
func (w *Table) PositionLongReadout() string {
	ymin, ymax := w.GetVisibleBoundaries()
	return fmt.Sprintf("%d,%d-%d/%d", w.list.Cursor()+1, ymin+1, ymax+1, w.list.Len())
}

// PositionShortReadout returns a percentage indicator on how far the songlist is scrolled.
// FIXME: move this into a positionreadout fragment
func (w *Table) PositionShortReadout() string {
	ymin, ymax := w.GetVisibleBoundaries()
	if ymin == 0 && ymax+1 == w.list.Len() {
		return `All`
	}
	if ymin == 0 {
		return `Top`
	}
	if ymax+1 == w.list.Len() {
		return `Bot`
	}
	fraction := float64(float64(ymin) / float64(w.list.Len()))
	percent := int(math.Floor(fraction * 100))
	return fmt.Sprintf("%2d%%", percent)
}

// SetColumns sets which columns that should be visible, and adjusts the sizes so they
// fit as close as possible to the median size of the content displayed.
func (w *Table) SetColumns(tags []string) {
	totalWidth, _ := w.Size()
	usedWidth := 0

	cols := w.list.Columns(tags)
	w.columns = make([]column, len(tags))

	for i, key := range tags {
		w.columns[i].col = cols[i]
		w.columns[i].key = key
		w.columns[i].width = cols[i].Median()
		usedWidth += w.columns[i].width
	}

	w.visibleColumns = tags

	if len(tags) == 0 {
		return
	}

	// right-hand column must have some space for readability
	w.columns[len(tags)-1].rightPadding = 1

	// log.Debugf("expanding column widths from %d to %d", usedWidth, totalWidth)

	// expand to size
	poolSize := len(tags)
	saturated := make([]bool, poolSize)

	// expand as long as there is space left
	for {
		for i := range tags {
			if usedWidth > totalWidth {
				return
			}
			if poolSize > 0 && saturated[i] {
				continue
			}
			col := w.columns[i]
			if poolSize > 0 && col.width > col.col.Max() {
				// log.Debugf("saturating column %s at width %d", tags[i], col.width)
				saturated[i] = true
				poolSize--
				continue
			}
			w.columns[i].width++
			// log.Debugf("increase column %s to width %d", tags[i], w.columns[i].width)
			usedWidth++
		}
	}
}

// ScrollViewport scrolls the viewport by delta rows, as far as possible.
// If movecursor is false, the cursor is kept pointing at the same song where
// possible. If true, the cursor is moved delta rows.
func (w *Table) ScrollViewport(delta int, movecursor bool) {
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
		w.list.MoveCursor(delta)
	}

	w.validateCursor()
}

// validateCursor ensures the cursor is within the allowable area without moving
// the viewport.
func (w *Table) validateCursor() {
	ymin, ymax := w.GetVisibleBoundaries()
	list := w.list
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
