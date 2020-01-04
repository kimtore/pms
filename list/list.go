package list

import (
	"fmt"
	"sort"
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
	CursorRow() Row
	MoveCursor(int)
	SetCursor(int)
	ValidateCursor(int, int)
}

type Metadata interface {
	ColumnNames() []string
	Columns([]string) []Column
	Name() string
	SetColumnNames([]string)
	SetName(string)
}

type List interface {
	Cursor
	Metadata
	Selectable
	Add(Row)
	Clear()
	InRange(int) bool
	Len() int
	Lock()
	NextOf([]string, int, int) int
	Row(int) Row
	RowNum(string) (int, error)
	SetUpdated()
	Sort([]string) error
	Unlock()
	Updated() time.Time
}

type Base struct {
	columnNames     []string
	columns         map[string]*Column
	cursor          int
	mutex           sync.Mutex
	name            string
	rows            []Row
	selection       map[int]struct{}
	sortKey         string
	updated         time.Time
	visualSelection [3]int
}

func New() *Base {
	s := &Base{}
	s.Clear()
	return s
}

func (s *Base) Clear() {
	s.rows = make([]Row, 0)
	s.columnNames = make([]string, 0)
	s.columns = make(map[string]*Column)
	s.ClearSelection()
}

func (s *Base) SetColumnNames(names []string) {
	s.columnNames = names
}

func (s *Base) ColumnNames() []string {
	names := make([]string, 0, len(s.columns))
	for key := range s.columns {
		names = append(names, key)
	}
	return names
}

func (s *Base) Columns(names []string) []Column {
	cols := make([]Column, len(names))
	for i, name := range names {
		if col, ok := s.columns[name]; ok {
			cols[i] = *col
		}
	}
	return cols
}

func (s *Base) Add(row Row) {
	s.rows = append(s.rows, row)
	for k, v := range row {
		if s.columns[k] == nil {
			s.columns[k] = &Column{}
		}
		s.columns[k].Add(v)
	}
}

func (s *Base) Row(n int) Row {
	if !s.InRange(n) {
		return nil
	}
	return s.rows[n]
}

func (s *Base) RowNum(id string) (int, error) {
	for n, row := range s.rows {
		if row.ID() == id {
			return n, nil
		}
	}
	return 0, fmt.Errorf("not found")
}

func (s *Base) Len() int {
	return len(s.rows)
}

// Implements sort.Interface
func (s *Base) Less(i, j int) bool {
	return s.rows[i][s.sortKey] < s.rows[j][s.sortKey]
}

// Implements sort.Interface
func (s *Base) Swap(i, j int) {
	row := s.rows[i]
	s.rows[i] = s.rows[j]
	s.rows[j] = row
}

// Sort first sorts unstable, then stable, by all columns provided.
// Retains cursor position.
func (s *Base) Sort(cols []string) error {
	if s.Len() < 2 {
		return nil
	}

	// Obtain row under cursor
	cursorRow := s.CursorRow()

	fn := sort.Sort
	for _, key := range cols {
		s.sortKey = key
		fn(s)
		fn = sort.Stable
	}

	// Restore cursor position to row previously selected
	rowNum, err := s.RowNum(cursorRow.ID())
	if err != nil {
		// panics here because the row with this id must also be found in the sorted list,
		// otherwise this is a bug.
		panic(err)
	}

	s.SetCursor(rowNum)

	return nil
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

// NextOf searches forwards or backwards for rows having different values in the specified tags.
// The index of the next song is returned.
func (s *Base) NextOf(tags []string, index int, direction int) int {
	offset := func(i int) int {
		if direction > 0 || i == 0 {
			return 0
		}
		return 1
	}

	ln := s.Len()
	index -= offset(index)
	row := s.Row(index)

LOOP:
	for ; index < ln && index >= 0; index += direction {
		for _, tag := range tags {
			if row[tag] != s.rows[index][tag] {
				break LOOP
			}
		}
	}

	return index + offset(index)
}

func (s *Base) Remove(index int) error {
	row := s.Row(index)
	if row == nil {
		return fmt.Errorf("out of bounds")
	}

	for k, v := range row {
		s.columns[k].Remove(v)
	}

	if index+1 == s.Len() {
		s.rows = s.rows[:index]
	} else {
		s.rows = append(s.rows[:index], s.rows[index+1:]...)
	}

	return nil
}

// RemoveIndices removes a selection of songs from the songlist, having the
// index defined by the int slice parameter.
func (s *Base) RemoveIndices(indices []int) error {
	// Ensure that indices are removed in reverse order
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))
	for _, i := range indices {
		if err := s.Remove(i); err != nil {
			return err
		}
	}
	return nil
}
