package topbar

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/utils"
	"github.com/gdamore/tcell/views"
)

type Matrix [][]Piece

// Fragment is the smallest possible unit in a topbar.
type Fragment interface {
	Draw(x, y int) int
	SetDirty(bool)
	SetView(views.View)
	Width() int
	style.Stylable
}

// Piece is one unit in the matrix.
type Piece struct {
	align     int
	fragments []Fragment
	padding   int
	view      views.View
	width     int
	style.Styled
}

// Contents of Pieces may be aligned to left, center or right.
const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)

func NewPiece() Piece {
	return Piece{
		padding:   1,
		fragments: make([]Fragment, 0),
	}
}

// MakePieces creates a new two-dimensional array of Piece objects.
func NewMatrix(width, height int) Matrix {
	matrix := make(Matrix, height)
	for y := 0; y < height; y++ {
		matrix[y] = make([]Piece, width)
	}
	return matrix
}

// autoAlign returns a best-guess align for a Piece.
func autoAlign(x, width int) int {
	switch x {
	case 0:
		return AlignLeft
	case width - 1:
		return AlignRight
	default:
		return AlignCenter
	}
}

// Expand ensures that all rows of a Matrix have the same length, and sets auto-alignment.
func (matrix Matrix) Expand() {
	width := 0
	for y := 0; y < len(matrix); y++ {
		width = utils.Max(width, len(matrix[y]))
	}
	for y := 0; y < len(matrix); y++ {
		for x := 0; x < width; x++ {
			if x >= len(matrix[y]) {
				p := NewPiece()
				matrix[y] = append(matrix[y], p)
			}
			matrix[y][x].SetAlign(autoAlign(x, width))
		}
	}
}

// Size returns the dimensions of a Matrix.
func (matrix Matrix) Size() (x, y int) {
	y = len(matrix)
	if y > 0 {
		x = len(matrix[0])
	}
	return
}

func (matrix Matrix) SetView(v views.View) {
	xmax, ymax := matrix.Size()
	for y := 0; y < ymax; y++ {
		for x := 0; x < xmax; x++ {
			matrix[y][x].SetView(v)
		}
	}
}

// Parse creates a new two-dimensional array of Piece objects based on lexer input.
func Parse(input string) (Matrix, error) {

	piece := NewPiece()
	matrix := NewMatrix(0, 0)
	row := make([]Piece, 0)
	pos := 0
	x := 0

	variable := false

	addPiece := func() {
		row = append(row, piece)
		piece = NewPiece()
		x += 1
	}

	addRow := func() {
		addPiece()
		matrix = append(matrix, row)
		row = make([]Piece, 0)
		x = 0
	}

	for {
		t, npos := lexer.NextToken(input[pos:])
		pos += npos
		s := t.String()

		if variable && t.Class != lexer.TokenIdentifier {
			return nil, fmt.Errorf("Unexpected '%s', expected identifier")
		}

		switch t.Class {
		case lexer.TokenEnd:
			goto end
		case lexer.TokenIdentifier:
			if variable {
				// FIXME: look up correct fragment
				piece.AddFragment(NewText(s))
				variable = false
			} else {
				piece.AddFragment(NewText(s))
			}
		case lexer.TokenSeparator:
			addPiece()
		case lexer.TokenStop:
			addRow()
		case lexer.TokenVariable:
			variable = true
			continue
		default:
			return nil, fmt.Errorf("Unexpected '%s', expected variable or identifier", s)
		}
	}

end:
	addRow()
	matrix.Expand()

	return matrix, nil
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

func (p *Piece) SetAlign(align int) {
	p.align = align
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
