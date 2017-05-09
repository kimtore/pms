package widgets

import (
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/topbar"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// Topbar is a widget that can display a variety of information, such as the
// currently playing song. It is composed of several pieces to form a
// two-dimensional matrix.
type Topbar struct {
	matrix topbar.Matrix
	width  int // width is the matrix width, as opposed to character width
	height int // height is both physical and matrix height

	view views.View
	style.Styled
	views.WidgetWatchers
}

// NewTopbar creates a new Topbar widget in the desired dimensions.
func NewTopbar() *Topbar {
	return &Topbar{
		width:  0,
		height: 0,
		matrix: topbar.NewMatrix(0, 0),
	}
}

// pieceWidth returns the auto-adjusted width of any topbar Piece.
func (w *Topbar) pieceWidth() int {
	if w.width == 0 {
		return 0
	}
	xmax, _ := w.Size()
	return xmax / w.width
}

// Setup sets up the topbar using the provided configuration string.
func (w *Topbar) SetMatrix(matrix topbar.Matrix) {
	matrix.SetView(w.view)
	matrix.SetStylesheet(w.Stylesheet())
	w.matrix = matrix
	w.width, w.height = w.matrix.Size()
	console.Log("Setting up new topbar with dimensions (%d, %d)", w.width, w.height)
}

// Draw draws all the pieces in the matrix, from top to bottom, right to left.
func (w *Topbar) Draw() {
	pieceWidth := w.pieceWidth()

	w.view.Fill(' ', w.Style("topbar"))

	for y := 0; y < w.height; y++ {
		for x := w.width - 1; x >= 0; x-- {
			w.matrix[y][x].Draw(x*pieceWidth, y, pieceWidth)
		}
	}
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
	w.matrix.SetView(w.view)
}
