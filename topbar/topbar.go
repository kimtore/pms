package topbar

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/utils"
	"github.com/gdamore/tcell/views"
)

type Matrix [][]*Piece

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

var fragments = map[string]func(api.API) Fragment{
	"artist":    NewArtist,
	"shortname": NewShortname,
	"version":   NewVersion,
}

func NewRow() []*Piece {
	return make([]*Piece, 0)
}

func NewPiece() *Piece {
	return &Piece{
		fragments: make([]Fragment, 0),
	}
}

// EmptyMatrix creates a new two-dimensional array of Piece objects.
func EmptyMatrix() Matrix {
	return make(Matrix, 0)
}

// NewMatrix creates a matrix based on a parsed topbar matrix statement.
func NewMatrix(matrixStmt *MatrixStatement, a api.API) (Matrix, error) {

	var frag Fragment
	matrix := EmptyMatrix()

	for _, rowStmt := range matrixStmt.Rows {
		row := NewRow()

		for _, pieceStmt := range rowStmt.Pieces {
			piece := NewPiece()

			for _, fragmentStmt := range pieceStmt.Fragments {

				if len(fragmentStmt.Variable) > 0 {
					ctor, ok := fragments[fragmentStmt.Variable]
					if !ok {
						return nil, fmt.Errorf("Unrecognized variable '${%s}'", fragmentStmt.Variable)
					}
					frag = ctor(a)
				} else {
					frag = NewText(fragmentStmt.Literal)
				}
				piece.AddFragment(frag)
			}
			row = append(row, piece)
		}
		matrix = append(matrix, row)
	}

	matrix.Expand()
	return matrix, nil
}

func Parse(a api.API, input string) (Matrix, error) {
	reader := strings.NewReader(input)
	parser := NewParser(reader)

	stmt, err := parser.ParseMatrix()
	if err != nil {
		return nil, err
	}

	return NewMatrix(stmt, a)
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

// mapPiece runs a callback function once for each piece in the matrix.
func (matrix Matrix) mapPiece(f func(*Piece)) {
	xmax, ymax := matrix.Size()
	for y := 0; y < ymax; y++ {
		for x := 0; x < xmax; x++ {
			f(matrix[y][x])
		}
	}
}

func (matrix Matrix) SetStylesheet(stylesheet style.Stylesheet) {
	matrix.mapPiece(func(p *Piece) {
		p.SetStylesheet(stylesheet)
	})
}

func (matrix Matrix) SetView(v views.View) {
	matrix.mapPiece(func(p *Piece) {
		p.SetView(v)
	})
}

// Draw draws all fragments.
func (p *Piece) Draw(x, y, width int) {
	//console.Log("Align %d says draw at xorig=%d, xalign=%d, width=%d, textwidth=%d", p.align, x, p.alignX(x, width), width, p.Width())
	x = p.alignX(x, width)
	for _, fragment := range p.fragments {
		x = fragment.Draw(x, y)
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
	return width
}

func (p *Piece) SetView(v views.View) {
	for _, fragment := range p.fragments {
		fragment.SetView(v)
	}
}

func (p *Piece) SetAlign(align int) {
	p.align = align
}

func (p *Piece) SetStylesheet(stylesheet style.Stylesheet) {
	p.Styled.SetStylesheet(stylesheet)
	for _, fragment := range p.fragments {
		fragment.SetStylesheet(stylesheet)
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
