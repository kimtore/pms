package widgets

import (
	"strings"

	"github.com/gdamore/tcell/views"
)

type ColumnheadersWidget struct {
	columns []column
	view    views.View

	widget
}

func NewColumnheadersWidget() (c *ColumnheadersWidget) {
	c = &ColumnheadersWidget{}
	c.columns = make([]column, 0)
	return
}

func (c *ColumnheadersWidget) SetColumns(cols []column) {
	c.columns = cols
}

func (c *ColumnheadersWidget) Draw() {
	x := 0
	y := 0
	for i := range c.columns {
		col := &c.columns[i]
		title := []rune(strings.Title(col.Tag))
		for p, r := range title {
			c.view.SetContent(x+p, y, r, nil, c.Style("header"))
		}
		x += col.Width()
	}
}

func (c *ColumnheadersWidget) SetView(v views.View) {
	c.view = v
}

func (c *ColumnheadersWidget) Size() (int, int) {
	x, y := c.view.Size()
	y = 1
	return x, y
}
