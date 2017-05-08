package widgets

import (
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/topbar"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// Topbar is a widget that can display a variety of information, such as the
// currently playing song. It is composed of several pieces to form a
// two-dimensional matrix.
type Topbar struct {
	pieces [][]topbar.Piece
	width  int // width is the matrix width, as opposed to character width
	height int // height is both physical and matrix height

	view views.View
	style.Styled
	views.WidgetWatchers
}

// NewTopbar creates a new Topbar widget in the desired dimensions.
func NewTopbar(width, height int) *Topbar {
	return &Topbar{
		width:  width,
		height: height,
		pieces: makePieces(width, height),
	}
}

// makePieces creates a new two-dimensional array of Piece objects.
func makePieces(width, height int) [][]topbar.Piece {
	pieces := make([][]topbar.Piece, height)
	for y := 0; y < height; y++ {
		pieces[y] = make([]topbar.Piece, width)
	}
	return pieces
}

// pieceWidth returns the auto-adjusted width of any topbar Piece.
func (w *Topbar) pieceWidth() int {
	xmax, _ := w.Size()
	return xmax / w.width
}

// FIXME: configurability
// Setup sets up a sample topbar.
func (w *Topbar) Setup() {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			p := topbar.NewPiece(topbar.AlignLeft + x)
			w.SetPiece(x, y, p)
		}
	}

	w.pieces[0][0].AddFragment(&topbar.Shortname{})
	w.pieces[0][0].AddFragment(&topbar.Version{})
}

// Draw draws all the pieces in the matrix, from top to bottom, right to left.
func (w *Topbar) Draw() {
	pieceWidth := w.pieceWidth()

	w.view.Fill(' ', w.Style("topbar"))

	for y := 0; y < w.height; y++ {
		for x := w.width - 1; x >= 0; x-- {
			w.pieces[y][x].Draw(x*pieceWidth, y, pieceWidth)
		}
	}
}

// SetPiece specifies that the given Piece should be drawn at the given matrix coordinates.
func (w *Topbar) SetPiece(x, y int, p topbar.Piece) {
	p.SetStylesheet(w.Stylesheet())
	w.pieces[y][x] = p
}

func (w *Topbar) HandleEvent(ev tcell.Event) bool {
	return false
}

func (w *Topbar) Size() (int, int) {
	x, _ := w.view.Size()
	return x, w.height
}

func (w *Topbar) Resize() {
}

func (w *Topbar) SetView(v views.View) {
	w.view = v
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			w.pieces[y][x].SetView(w.view)
		}
	}
}
