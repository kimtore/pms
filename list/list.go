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
		cols[i] = *s.columns[name]
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

func (s *Base) Sort(cols []string) error {
	fn := sort.Sort
	for _, key := range cols {
		s.sortKey = key
		fn(s)
		fn = sort.Stable
	}
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
