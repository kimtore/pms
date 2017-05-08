package topbar

import (
	"github.com/ambientsound/pms/style"
	"github.com/gdamore/tcell/views"
)

// Fragment is the smallest possible unit in a topbar.
type Fragment interface {
	Draw(x, y int) int
	SetDirty(bool)
	SetView(views.View)
	Width() int
	style.Stylable
}

// Contents of Pieces may be aligned to left, center or right.
const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)

// Piece is one unit in the matrix.
type Piece struct {
	align     int
	fragments []Fragment
	padding   int
	view      views.View
	width     int
	style.Styled
}

func NewPiece(align int) Piece {
	return Piece{
		align:     align,
		padding:   1,
		fragments: make([]Fragment, 0),
	}
}

// Draw draws all fragments.
func (p *Piece) Draw(x, y, width int) {
	//console.Log("Align %d says draw at xorig=%d, xalign=%d, width=%d, textwidth=%d", p.align, x, p.alignX(x, width), width, p.Width())
	x = p.alignX(x, width)
	for _, fragment := range p.fragments {
		x = fragment.Draw(x, y) + p.padding
	}
}

// AddFragment appends the given fragment to the list of fragments that should be drawn in this piece.
func (p *Piece) AddFragment(f Fragment) {
	f.SetStylesheet(p.Stylesheet())
	f.SetView(p.view)
	p.fragments = append(p.fragments, f)
}

// Width returns the total text width of all fragments.
func (p *Piece) Width() int {
	width := 0
	for _, fragment := range p.fragments {
		width += fragment.Width()
	}
	width += p.padWidth()
	return width
}

// padWidth returns the total whitespace length between fragments.
func (p *Piece) padWidth() int {
	if len(p.fragments) == 0 {
		return 0
	}
	return p.padding * (len(p.fragments) - 1)
}

func (p *Piece) SetView(v views.View) {
	for _, fragment := range p.fragments {
		fragment.SetView(v)
	}
}

// drawX returns the draw start position.
func (p *Piece) alignX(x, width int) int {
	switch p.align {
	case AlignLeft:
		return x
	case AlignCenter:
		return x + (width / 2) - (p.Width() / 2)
	case AlignRight:
		return x + width - p.Width()
	default:
		return x
	}
}
