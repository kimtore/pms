package term

import (
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
		y2:     y + height,
		width:  width,
		height: height,
	}
}

// Fill the entire canvas with a character.
func (c Canvas) Fill(r rune, s Style) {
	fg, bg := s.Attr()
	for y := c.y1; y <= c.y2; y++ {
		for x := c.x1; x <= c.x2; x++ {
			termbox.SetCell(x, y, r, fg, bg)
		}
	}
}

// Print text onto the canvas.
func (c Canvas) SetCell(x, y int, r rune, st Style) int {
	x += c.x1
	y += c.y1
	fg, bg := st.Attr()
	termbox.SetCell(x, y, r, fg, bg)
	return x + 1
}

// Print text onto the canvas.
func (c Canvas) Print(x, y int, s string, st Style) int {
	x += c.x1
	y += c.y1
	fg, bg := st.Attr()
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
	return c.Width(), c.Height()
}

// Return the canvas height.
func (c Canvas) Height() int {
	return c.y2 - c.y1
}

// Return the canvas width.
func (c Canvas) Width() int {
	return c.x2 - c.x1
}
