package widgets

import (
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/utils"
)

type column struct {
	Tag        string
	items      int
	totalWidth int
	maxWidth   int
	mid        int
	width      int
}

type columns []column

func (c *column) Set(s songlist.Songlist) {
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

func (c *column) Mid() int {
	return c.mid
}

func (c *column) MaxWidth() int {
	return c.maxWidth
}

func (c *column) Width() int {
	return c.width
}

func (c *column) SetWidth(width int) {
	c.width = width
}
