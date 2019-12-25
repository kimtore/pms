package list

import (
	"sync"
	"time"
)

type Selectable interface {
	ClearSelection()
	CommitVisualSelection()
	DisableVisualSelection()
	EnableVisualSelection()
	HasVisualSelection() bool
	Selected(int) bool
	SelectionIndices() []int
	SetSelected(int, bool)
	SetVisualSelection(int, int, int)
	ToggleVisualSelection()
}

type Cursor interface {
	Cursor() int
	MoveCursor(int)
	SetCursor(int)
	ValidateCursor(int, int)
}

type Metadata interface {
	Name() string
	SetName(string) error
	Columns() []Column
}

type List interface {
	Selectable
	Cursor
	InRange(int) bool
	Len() int
	Lock()
	SetUpdated()
	Sort([]string) error
	Unlock()
	Updated() time.Time
}

type Item interface {
	Len() int
}

type Base struct {
	items           []Item
	columns         []Column
	cursor          int
	mutex           sync.Mutex
	name            string
	selection       map[int]struct{}
	updated         time.Time
	visualSelection [3]int
}

func (s *Base) Len() int {
	return len(s.items)
}

// InRange returns true if the provided index is within list range, false otherwise.
func (s *Base) InRange(index int) bool {
	return index >= 0 && index < s.Len()
}

func (s *Base) Lock() {
	s.mutex.Lock()
}

func (s *Base) Unlock() {
	s.mutex.Unlock()
}

func (s *Base) Name() string {
	return s.name
}

func (s *Base) SetName(name string) {
	s.name = name
}

// Updated returns the timestamp of when this songlist was last updated.
func (s *Base) Updated() time.Time {
	return s.updated
}

// SetUpdated sets the update timestamp of the songlist.
func (s *Base) SetUpdated() {
	s.updated = time.Now()
}
