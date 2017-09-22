package term

import (
	"github.com/gdamore/tcell"
	termbox "github.com/nsf/termbox-go"
)

// Canvas draws text on the terminal. It s represented as a rectangle,
// providing a simple API on top.
type Canvas struct {
	x1     int
	x2     int
	y1     int
	y2     int
	width  int
	height int
}

// NewCanvas returns Canvas.
func NewCanvas(x, y, width, height int) Canvas {
	return Canvas{
		x1:     x,
		y1:     y,
		x2:     x + width,
		y2:     y + width,
		width:  width,
		height: height,
	}
}

// Fill the entire canvas with a character.
func (c Canvas) Fill(r rune, s tcell.Style) {
	var fg, bg termbox.Attribute
	for y := c.y1; y <= c.y2; y++ {
		for x := c.x1; x <= c.x2; x++ {
			termbox.SetCell(x, y, r, fg, bg)
			x++
		}
		y++
	}
}

// Print text onto the canvas.
func (c Canvas) Print(x, y int, s string, st tcell.Style) int {
	var fg, bg termbox.Attribute
	x += c.x1
	y += c.y1
	//fg, bg, _ := s.Decompose()
	for _, ch := range s {
		termbox.SetCell(x, y, ch, fg, bg)
		x++
		if x > c.x2 {
			break
		}
	}
	return x
}

// Return the canvas size.
func (c Canvas) Size() (width, height int) {
	width = c.x2 - c.x1
	height = c.y2 - c.y1
	return
}
