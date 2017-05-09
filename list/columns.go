package list

import (
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/utils"
)

type Column struct {
	Tag        string
	items      int
	totalWidth int
	maxWidth   int
	mid        int
	width      int
}

type Columns []Column

func (c *Column) Set(s songlist.Songlist) {
	for _, song := range s.Songs() {
		l := len(song.Tags[c.Tag])
		if l > 0 {
			c.items++
			c.totalWidth += l
			c.maxWidth = utils.Max(c.maxWidth, l)
		}
	}
	if c.items == 0 {
		c.mid = 0
	} else {
		c.mid = c.totalWidth / c.items
	}
}

func (c *Column) Mid() int {
	return c.mid
}

func (c *Column) MaxWidth() int {
	return c.maxWidth
}

func (c *Column) Width() int {
	return c.width
}

func (c *Column) SetWidth(width int) {
	c.width = width
}
