package list

import (
	"sort"

	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
)

// List maintains cursor, columns and selection state for a Songlist.
// TODO: merge into Songlist itself?
type List struct {
	Columns         Columns
	Cursor          int
	Selection       map[int]struct{}
	Songlist        songlist.Songlist
	visualSelection [3]int
	Ymax            int
	Ymin            int
}

func New(s songlist.Songlist) *List {
	l := &List{}
	l.Songlist = s
	l.Columns = make(Columns, 0)
	l.ClearSelection()
	return l
}

// CursorSong returns the song currently selected by the cursor.
func (w *List) CursorSong() *song.Song {
	return w.Songlist.Song(w.Cursor)
}

// ManuallySelected returns true if the given song index is selected through manual selection.
func (w *List) ManuallySelected(i int) bool {
	_, ok := w.Selection[i]
	return ok
}

// VisuallySelected returns true if the given song index is selected through visual selection.
func (w *List) VisuallySelected(i int) bool {
	return w.visualSelection[0] <= i && i <= w.visualSelection[1]
}

// Selected returns true if the given song index is selected, either through
// visual selection or manual selection. If the song is doubly selected, the
// selection is inversed.
func (w *List) Selected(i int) bool {
	a := w.ManuallySelected(i)
	b := w.VisuallySelected(i)
	return (a || b) && a != b
}

// SelectionIndices returns a slice of ints holding the position of each
// element in the current selection. If no elements are selected, the cursor
// position is returned.
func (w *List) SelectionIndices() []int {
	selection := make([]int, 0, w.Songlist.Len())
	max := w.Songlist.Len()
	for i := 0; i < max; i++ {
		if w.Selected(i) {
			selection = append(selection, i)
		}
	}
	if len(selection) == 0 && w.Songlist.Len() > 0 {
		selection = append(selection, w.Cursor)
	}
	selection = sort.IntSlice(selection)
	return selection
}

// validateVisualSelection makes sure the visual selection stays in range of
// the songlist size.
func (w *List) validateVisualSelection(ymin, ymax, ystart int) (int, int, int) {
	if w.Songlist.Len() == 0 || ymin < 0 || ymax < 0 || !w.Songlist.InRange(ystart) {
		return -1, -1, -1
	}
	if !w.Songlist.InRange(ymin) {
		ymin = 0
	}
	if !w.Songlist.InRange(ymax) {
		ymax = w.Songlist.Len() - 1
	}
	return ymin, ymax, ystart
}

// VisualSelection returns the min, max, and start position of visual select.
func (w *List) VisualSelection() (int, int, int) {
	return w.visualSelection[0], w.visualSelection[1], w.visualSelection[2]
}

// SetVisualSelection sets the range of the visual selection. Use negative
// integers to un-select all visually selected songs.
func (w *List) SetVisualSelection(ymin, ymax, ystart int) {
	w.visualSelection[0] = ymin
	w.visualSelection[1] = ymax
	w.visualSelection[2] = ystart
}

// HasVisualSelection returns true if the songlist is in visual selection mode.
func (w *List) HasVisualSelection() bool {
	return w.visualSelection[0] >= 0 && w.visualSelection[1] >= 0
}

// SetSelection sets the selected status of a single song.
func (w *List) SetSelected(i int, selected bool) {
	var x struct{}
	_, ok := w.Selection[i]
	if ok == selected {
		return
	}
	if selected {
		w.Selection[i] = x
	} else {
		delete(w.Selection, i)
	}
}

// CommitVisualSelection converts the visual selection to manual selection.
func (w *List) CommitVisualSelection() {
	if !w.HasVisualSelection() {
		return
	}
	for key := w.visualSelection[0]; key <= w.visualSelection[1]; key++ {
		selected := w.Selected(key)
		w.SetSelected(key, selected)
	}
}

// ClearSelection removes all selection.
func (w *List) ClearSelection() {
	w.Selection = make(map[int]struct{}, 0)
	w.visualSelection = [3]int{-1, -1, -1}
}
