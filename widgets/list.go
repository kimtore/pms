package widgets

import (
	"sort"

	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
)

type list struct {
	columns         columns
	cursor          int
	selection       map[int]struct{}
	songlist        songlist.Songlist
	visualSelection [3]int
	ymax            int
	ymin            int
}

func newList(s songlist.Songlist) *list {
	l := &list{}
	l.songlist = s
	l.columns = make(columns, 0)
	l.ClearSelection()
	return l
}

// CursorSong returns the song currently selected by the cursor.
func (w *list) CursorSong() *song.Song {
	return w.songlist.Song(w.cursor)
}

// ManuallySelected returns true if the given song index is selected through manual selection.
func (w *list) ManuallySelected(i int) bool {
	_, ok := w.selection[i]
	return ok
}

// VisuallySelected returns true if the given song index is selected through visual selection.
func (w *list) VisuallySelected(i int) bool {
	return w.visualSelection[0] <= i && i <= w.visualSelection[1]
}

// Selected returns true if the given song index is selected, either through
// visual selection or manual selection. If the song is doubly selected, the
// selection is inversed.
func (w *list) Selected(i int) bool {
	a := w.ManuallySelected(i)
	b := w.VisuallySelected(i)
	return (a || b) && a != b
}

// SelectionIndices returns a slice of ints holding the position of each
// element in the current selection. If no elements are selected, the cursor
// position is returned.
func (w *list) SelectionIndices() []int {
	selection := make([]int, 0, w.songlist.Len())
	max := w.songlist.Len()
	for i := 0; i < max; i++ {
		if w.Selected(i) {
			selection = append(selection, i)
		}
	}
	if len(selection) == 0 && w.songlist.Len() > 0 {
		selection = append(selection, w.cursor)
	}
	sort.Slice(selection, func(i, j int) bool { return selection[i] < selection[j] })
	return selection
}

// validateVisualSelection makes sure the visual selection stays in range of
// the songlist size.
func (w *list) validateVisualSelection(ymin, ymax, ystart int) (int, int, int) {
	if w.songlist.Len() == 0 || ymin < 0 || ymax < 0 || !w.songlist.InRange(ystart) {
		return -1, -1, -1
	}
	if !w.songlist.InRange(ymin) {
		ymin = 0
	}
	if !w.songlist.InRange(ymax) {
		ymax = w.songlist.Len() - 1
	}
	return ymin, ymax, ystart
}

// VisualSelection returns the min, max, and start position of visual select.
func (w *list) VisualSelection() (int, int, int) {
	return w.visualSelection[0], w.visualSelection[1], w.visualSelection[2]
}

// SetVisualSelection sets the range of the visual selection. Use negative
// integers to un-select all visually selected songs.
func (w *list) SetVisualSelection(ymin, ymax, ystart int) {
	w.visualSelection[0] = ymin
	w.visualSelection[1] = ymax
	w.visualSelection[2] = ystart
}

// HasVisualSelection returns true if the songlist is in visual selection mode.
func (w *list) HasVisualSelection() bool {
	return w.visualSelection[0] >= 0 && w.visualSelection[1] >= 0
}

// SetSelection sets the selected status of a single song.
func (w *list) SetSelected(i int, selected bool) {
	var x struct{}
	_, ok := w.selection[i]
	if ok == selected {
		return
	}
	if selected {
		w.selection[i] = x
	} else {
		delete(w.selection, i)
	}
}

func (w *list) ClearSelection() {
	w.selection = make(map[int]struct{}, 0)
	w.visualSelection = [3]int{-1, -1, -1}
}
