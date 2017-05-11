package songlist

import (
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/utils"
)

type Column struct {
	Tag        string
	items      int
	totalWidth int
	maxWidth   int
	avg        int
	width      int
}

type Columns []Column

type ColumnMap map[string]*Column

// Set calculates
func (c *Column) Set(s Songlist) {
	c.Reset()
	for _, song := range s.Songs() {
		c.Add(song)
	}
}

// Add a single song's width to the total and maximum width.
func (c *Column) Add(song *song.Song) {
	l := len(song.Tags[c.Tag])
	if l == 0 {
		return
	}
	c.avg = 0
	c.items++
	c.totalWidth += l
	c.maxWidth = utils.Max(c.maxWidth, l)
}

// Remove a single song's tag width from the total and maximum width.
func (c *Column) Remove(song *song.Song) {
	l := len(song.Tags[c.Tag])
	if l == 0 {
		return
	}
	c.avg = 0
	c.items--
	c.totalWidth -= l
	// FIXME: c.maxWidth is not updated
}

// Reset sets all values to zero.
func (c *Column) Reset() {
	c.items = 0
	c.totalWidth = 0
	c.maxWidth = 0
	c.avg = 0
	c.width = 0
}

// Weight returns the relative usefulness of this column. It might happen that
// a tag appears rarely, but is very long. In this case we reduce the field so
// that other tags get more space.
func (c *Column) Weight(max int) float64 {
	return float64(c.items) / float64(max)
}

// Avg returns the average length of the tag values in this column.
func (c *Column) Avg() int {
	if c.avg == 0 {
		if c.items == 0 {
			c.avg = 0
		} else {
			c.avg = c.totalWidth / c.items
		}
	}
	return c.avg
}

// MaxWidth returns the length of the longest tag value in this column.
func (c *Column) MaxWidth() int {
	return c.maxWidth
}

// Width returns the column width.
func (c *Column) Width() int {
	return c.width
}

// SetWidth sets the width that the column should consume.
func (c *Column) SetWidth(width int) {
	c.width = width
}

// expand adjusts the column widths equally between the different columns,
// giving affinity to weight.
func (columns Columns) expand(totalWidth int) {
	if len(columns) == 0 {
		return
	}

	usedWidth := 0
	poolSize := len(columns)
	saturated := make([]bool, poolSize)

	// Start with the average value
	for i := range columns {
		avg := columns[i].Avg()
		columns[i].SetWidth(avg)
		usedWidth += avg
	}

	// expand as long as there is space left
	for {
		for i := range columns {
			if usedWidth > totalWidth {
				return
			}
			if poolSize > 0 && saturated[i] {
				continue
			}
			col := &columns[i]
			if poolSize > 0 && col.Width() > col.MaxWidth() {
				saturated[i] = true
				poolSize--
				continue
			}
			col.SetWidth(col.Width() + 1)
			usedWidth++
		}
	}
}

// Add adds song tags to all applicable columns.
func (c ColumnMap) Add(song *song.Song) {
	for tag := range song.StringTags {
		c[tag].Add(song)
	}
}

// Remove removes song tags from all applicable columns.
func (c ColumnMap) add(song *song.Song) {
	for tag := range song.StringTags {
		c[tag].Remove(song)
	}
}

// ensureColumns makes sure that all of a song's tags exists in the column map.
func (s *BaseSonglist) ensureColumns(song *song.Song) {
	for tag := range song.StringTags {
		if _, ok := s.columns[tag]; !ok {
			s.columns[tag] = &Column{}
		}
	}
}

// Columns returns a slice of columns, containing only the columns which has
// the specified tags.
func (s *BaseSonglist) Columns(columns []string) Columns {
	cols := make(Columns, 0)
	for _, tag := range columns {
		col := s.columns[tag]
		if col == nil {
			col = &Column{Tag: tag}
		}
		cols = append(cols, *col)
	}
	return cols
}
