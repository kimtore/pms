package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type ColumnheadersWidget struct {
	columns []column
	view    views.View

	views.WidgetWatchers
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
	style := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	x := 0
	y := 0
	for i := range c.columns {
		col := &c.columns[i]
		for p := 0; p < len(col.Tag); p++ {
			c.view.SetContent(x+p, y, rune(col.Tag[p]), nil, style)
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

func (c *ColumnheadersWidget) Resize() {
	// Handled by SongListWidget
}

func (c *ColumnheadersWidget) HandleEvent(ev tcell.Event) bool {
	return false
}
