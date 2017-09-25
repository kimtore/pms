package widgets

import (
	"strings"
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/term"
	"github.com/ambientsound/pms/utils"
)

// SonglistWidget is a tcell widget which draws a Songlist on the screen. It
// maintains a list of songlists which can be cycled through.
type SonglistWidget struct {
	api      api.API
	columns  songlist.Columns
	lastDraw time.Time
	viewy    int
	tags     []string

	canvas
	style.Styled
}

// NewSonglistWidget returns SonglistWidget.
func NewSonglistWidget(a api.API) (w *SonglistWidget) {
	return &SonglistWidget{
		api: a,
	}
}

// Name returns a human-readable name for this songlist.
func (w *SonglistWidget) Name() string {
	return w.List().Name()
}

// SetColumns sets which columns that should be visible
func (w *SonglistWidget) SetColumns(tags []string) {
	w.tags = tags
	w.Resize()
}

// Resize expands the tag column widths.
func (w *SonglistWidget) Resize() {
	w.columns = w.List().Columns(w.tags)
	w.columns.Expand(w.c.Width())
}

// Panel is a shorthand function, returning the panel assigned to this songlist
// widget (left or right).
func (w *SonglistWidget) Panel() *songlist.Collection {
	return w.api.Db().Panel()
}

// List returns the active panel's songlist.
func (w *SonglistWidget) List() songlist.Songlist {
	return w.Panel().Current()
}

func (w *SonglistWidget) drawNext(x, y, strmin, strmax int, runes []rune, style term.Style) int {
	strmin = utils.Min(len(runes), strmin)
	n := 0
	for n < strmin {
		w.c.SetCell(x, y, runes[n], style)
		n++
		x++
	}
	for n < strmax {
		w.c.SetCell(x, y, ' ', style)
		n++
		x++
	}
	return x
}

func (w *SonglistWidget) drawOneTagLine(x, y, xmax int, s *song.Song, tag string, defaultStyle string, style term.Style, lineStyled bool) int {
	if !lineStyled {
		style = w.Style(defaultStyle)
	}

	runes := s.Tags[tag]
	strmin := len(runes)

	return w.drawNext(x, y, strmin, xmax+1, runes, style)
}

func (w *SonglistWidget) Draw() {
	//console.Log("Draw() in songlist widget")
	list := w.List()
	if list == nil || list.Songs() == nil {
		console.Log("BUG: nil list, aborting draw!")
		return
	}

	// Blank screen first
	w.c.Fill(' ', w.Style("default"))

	// Check if the current panel's songlist has changed.
	if w.Panel().Updated().After(w.lastDraw) {
		tags := strings.Split(w.api.Options().StringValue("columns"), ",")
		//cols := list.Columns(tags)
		w.SetColumns(tags)
	} else if list.Updated().Before(w.lastDraw) {
		//console.Log("SonglistWidget::Draw(): not drawing, already drawn")
		//return
	}

	// Make sure viewport shows the cursor.
	w.viewportToCursor()

	// Update draw time
	w.lastDraw = time.Now()

	// Calculate boundaries
	xmax, ymax := w.c.Size()
	xmax += 1
	ymax = w.Bottom()

	currentSong := w.api.Song()
	style := w.Style("default")
	cursor := false

	y := 0

	// Loop through top of viewport to end of viewport or end of list, whichever comes first.
	for i := w.viewy; i <= ymax; i++ {

		lineStyled := true
		s := list.Song(i)
		if s == nil {
			// Sometimes happens under race conditions; just abort drawing
			console.Log("Attempting to draw nil song on index %d, aborting draw due to possible race condition.", i)
			return
		}

		// Style based on song's role
		cursor = i == list.Cursor()
		switch {
		case cursor:
			style = w.Style("cursor")
		case list.IndexAtSong(i, currentSong):
			style = w.Style("currentSong")
		case list.Selected(i):
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

		y++
	}
}

// Size returns the dimensions of the widget.
func (w *SonglistWidget) Size() (int, int) {
	return w.c.Size()
}

// Top returns the index from the dataset that is at the top of the viewport.
func (w *SonglistWidget) Top() int {
	return w.viewy
}

// Bottom returns the index from the dataset that is at the bottom of the
// viewport. If the dataset is empty, return -1.
func (w *SonglistWidget) Bottom() int {
	return utils.Min(
		w.viewy+w.c.Height()-1, // viewport end
		w.List().Len()-1,       // dataset size
	)
}

// Scroll scrolls the viewport by delta rows, as far as possible.
// If movecursor is false, the cursor is kept pointing at the same song where
// possible. If true, the cursor is moved delta rows.
func (w *SonglistWidget) Scroll(delta int, movecursor bool) {
	if movecursor {
		w.List().MoveCursor(delta)
	}
	w.SetViewport(w.viewy + delta)
	console.Log("Scroll: top=%d, bottom=%d, cursor=%d", w.Top(), w.Bottom(), w.List().Cursor())
	w.cursorToViewport()
	console.Log("Scroll: top=%d, bottom=%d, cursor=%d", w.Top(), w.Bottom(), w.List().Cursor())
}

// SetViewport sets the top position of the viewport.
func (w *SonglistWidget) SetViewport(pos int) {
	w.viewy = pos
	w.Validate(0, w.Bottom())
}

// Validate ensures that the viewport dimensions are within boundaries.
func (w *SonglistWidget) Validate(ymin, ymax int) {
	if w.Bottom() > ymax {
		w.viewy = ymax - w.c.Height()
	}
	if w.Top() < ymin {
		w.viewy = ymin
	}
}

// viewportToCursor moves the visible viewport so that the cursor is made visible.
// If the 'center' option is enabled, the viewport is centered on the cursor.
func (w *SonglistWidget) viewportToCursor() {
	cursor := w.List().Cursor()

	// Make the cursor visible.
	if !w.api.Options().BoolValue("center") {
		w.MakeVisible(cursor)
		return
	}

	// If 'center' is on, make the cursor centered.
	_, height := w.c.Size()
	half := height / 2
	min := utils.Max(0, cursor-half)
	max := utils.Min(w.List().Len()-1, cursor+half)
	w.MakeVisible(min)
	w.MakeVisible(max)
}

// cursorToViewport ensures the cursor is within the allowable area without moving
// the viewport.
func (w *SonglistWidget) cursorToViewport() {
	top := w.Top()
	bottom := w.Bottom()
	list := w.List()
	cursor := list.Cursor()

	if w.api.Options().BoolValue("center") {
		// When 'center' is on, move cursor to the centre of the viewport
		target := cursor
		lowerbound := (top + bottom) / 2
		upperbound := lowerbound
		if top <= 0 {
			// We are scrolled to the top, so the cursor is allowed to go above
			// the middle of the viewport
			lowerbound = 0
		}
		if bottom >= list.Len()-1 {
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
		if cursor < top {
			list.SetCursor(top)
		} else if cursor > bottom {
			list.SetCursor(bottom)
		}
	}
}

// MakeVisible adjusts the viewport by the minimal amount possible, to ensure
// that a certain index is visible.
func (w *SonglistWidget) MakeVisible(y int) {
	height := w.c.Height()
	if y < w.List().Len() && y >= w.viewy+height {
		w.viewy = y - (height - 1)
	}
	if y >= 0 && y < w.viewy {
		w.viewy = y
	}
}

// PositionReadout returns a combination of PositionLongReadout() and PositionShortReadout().
// FIXME: move this into a positionreadout fragment
func (w *SonglistWidget) PositionReadout() string {
	//return fmt.Sprintf("%s    %s", w.PositionLongReadout(), w.PositionShortReadout())
	return ""
}

// PositionLongReadout returns a formatted string containing the visible song
// range as well as the total number of songs.
// FIXME: move this into a positionreadout fragment
func (w *SonglistWidget) PositionLongReadout() string {
	//ymin, ymax := w.GetVisibleBoundaries()
	//return fmt.Sprintf("%d,%d-%d/%d", w.List().Cursor()+1, ymin+1, ymax+1, w.List().Len())
	return ""
}

// PositionShortReadout returns a percentage indicator on how far the songlist is scrolled.
// FIXME: move this into a positionreadout fragment
func (w *SonglistWidget) PositionShortReadout() string {
	/*
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
	*/
	return ""
}
