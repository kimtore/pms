package songlist

import (
	"sort"

	"github.com/ambientsound/pms/console"
)

// ManuallySelected returns true if the given song index is selected through manual selection.
func (s *BaseSonglist) ManuallySelected(i int) bool {
	_, ok := s.selection[i]
	return ok
}

// VisuallySelected returns true if the given song index is selected through visual selection.
func (s *BaseSonglist) VisuallySelected(i int) bool {
	return s.visualSelection[0] <= i && i <= s.visualSelection[1]
}

// Selected returns true if the given song index is selected, either through
// visual selection or manual selection. If the song is doubly selected, the
// selection is inversed.
func (s *BaseSonglist) Selected(i int) bool {
	a := s.ManuallySelected(i)
	b := s.VisuallySelected(i)
	return (a || b) && a != b
}

// SelectionIndices returns a slice of ints holding the position of each
// element in the current selection. If no elements are selected, the cursor
// position is returned.
func (s *BaseSonglist) SelectionIndices() []int {
	selection := make([]int, 0, s.Len())
	max := s.Len()
	for i := 0; i < max; i++ {
		if s.Selected(i) {
			selection = append(selection, i)
		}
	}
	if len(selection) == 0 && s.Len() > 0 {
		selection = append(selection, s.Cursor())
	}
	selection = sort.IntSlice(selection)
	return selection
}

// SetSelection sets the selected status of a single song.
func (s *BaseSonglist) SetSelected(i int, selected bool) {
	var x struct{}
	_, ok := s.selection[i]
	if ok == selected {
		return
	}
	if selected {
		s.selection[i] = x
	} else {
		delete(s.selection, i)
	}
}

// CommitVisualSelection converts the visual selection to manual selection.
func (s *BaseSonglist) CommitVisualSelection() {
	if !s.HasVisualSelection() {
		return
	}
	for key := s.visualSelection[0]; key <= s.visualSelection[1]; key++ {
		selected := s.Selected(key)
		s.SetSelected(key, selected)
	}
}

// ClearSelection removes all selection.
func (s *BaseSonglist) ClearSelection() {
	s.selection = make(map[int]struct{}, 0)
	s.visualSelection = [3]int{-1, -1, -1}
	// FIXME
	//PostEventModeSync(w, MultibarModeNormal)
}

// Selection returns the current selection as a new Songlist.
func (s *BaseSonglist) Selection() Songlist {
	indices := s.SelectionIndices()
	dest := New()
	for _, i := range indices {
		if song := s.Song(i); song != nil {
			dest.Add(song)
		} else {
			console.Log("SelectionIndices() returned an integer '%d' that resulted in a nil song, ignoring", i)
		}
	}
	return dest
}

// validateVisualSelection makes sure the visual selection stays in range of
// the songlist size.
func (s *BaseSonglist) validateVisualSelection(ymin, ymax, ystart int) (int, int, int) {
	if s.Len() == 0 || ymin < 0 || ymax < 0 || !s.InRange(ystart) {
		return -1, -1, -1
	}
	if !s.InRange(ymin) {
		ymin = 0
	}
	if !s.InRange(ymax) {
		ymax = s.Len() - 1
	}
	return ymin, ymax, ystart
}

// VisualSelection returns the min, max, and start position of visual select.
func (s *BaseSonglist) VisualSelection() (int, int, int) {
	return s.visualSelection[0], s.visualSelection[1], s.visualSelection[2]
}

// SetVisualSelection sets the range of the visual selection. Use negative
// integers to un-select all visually selected songs.
func (s *BaseSonglist) SetVisualSelection(ymin, ymax, ystart int) {
	s.visualSelection[0] = ymin
	s.visualSelection[1] = ymax
	s.visualSelection[2] = ystart
}

// HasVisualSelection returns true if the songlist is in visual selection mode.
func (s *BaseSonglist) HasVisualSelection() bool {
	return s.visualSelection[0] >= 0 && s.visualSelection[1] >= 0
}

// EnableVisualSelection sets start and stop of the visual selection to the
// cursor position.
func (s *BaseSonglist) EnableVisualSelection() {
	cursor := s.Cursor()
	s.SetVisualSelection(cursor, cursor, cursor)
}

// DisableVisualSelection disables visual selection.
func (s *BaseSonglist) DisableVisualSelection() {
	s.SetVisualSelection(-1, -1, -1)
}

// ToggleVisualSelection toggles visual selection on and off.
func (s *BaseSonglist) ToggleVisualSelection() {
	if !s.HasVisualSelection() {
		s.EnableVisualSelection()
	} else {
		s.DisableVisualSelection()
	}
}

// expandVisualSelection sets the visual selection boundaries from where it
// started to the current cursor position.
func (s *BaseSonglist) expandVisualSelection() {
	if !s.HasVisualSelection() {
		return
	}
	ymin, ymax, ystart := s.VisualSelection()
	switch {
	case s.Cursor() < ystart:
		ymin, ymax = s.Cursor(), ystart
	case s.Cursor() > ystart:
		ymin, ymax = ystart, s.Cursor()
	default:
		ymin, ymax = ystart, ystart
	}
	s.SetVisualSelection(ymin, ymax, ystart)
}
