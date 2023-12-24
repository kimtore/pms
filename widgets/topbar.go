package widgets

import (
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/topbar"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	`github.com/mattn/go-runewidth`
)

// Pieces may be aligned to left, center or right.
const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)

// Topbar is a widget that can display a variety of information, such as the
// currently playing song. It is composed of several pieces to form a
// two-dimensional matrix.
type Topbar struct {
	matrix *topbar.MatrixStatement
	height int // height is both physical and matrix height

	view views.View
	style.Styled
	views.WidgetWatchers
}

// NewTopbar creates a new Topbar widget in the desired dimensions.
func NewTopbar() *Topbar {
	return &Topbar{
		height: 0,
		matrix: &topbar.MatrixStatement{},
	}
}

// Setup sets up the topbar using the provided configuration string.
func (w *Topbar) SetMatrix(matrix *topbar.MatrixStatement) {
	w.matrix = matrix
	w.height = len(matrix.Rows)
	console.Log("Setting up new topbar with height %d", w.height)
}

// Draw draws all the pieces in the matrix, from top to bottom, right to left.
func (w *Topbar) Draw() {
	xmax, _ := w.Size()

	// Blank screen first
	w.view.Fill(' ', w.Style("topbar"))

	for y, rowStmt := range w.matrix.Rows {
		// Calculate window buffer width
		pieces := len(rowStmt.Pieces)
		if pieces == 0 {
			continue
		}

		for piece, pieceStmt := range rowStmt.Pieces {
			// Reset X position to start of window buffer, and align left,
			// center or right.
			align := autoAlign(piece, pieces)
			textWidth := pieceTextWidth(pieceStmt)
			x := getPiecesStartX(piece, pieces, xmax)
			x2 := getPiecesStartX(piece+1, pieces, xmax)
			x = alignX(x, x2-x, textWidth, align)

			for _, fragmentStmt := range pieceStmt.Fragments {
				frag := fragmentStmt.Instance
				text, styleStr := frag.Text()
				style := w.Style(styleStr)
				x = w.drawNext(x, y, text, style)
			}
		}
	}
}

// drawNext draws a string and returns the resulting X position.
func (w *Topbar) drawNext(x, y int, s string, style tcell.Style) int {
	for _, r := range s {
		w.view.SetContent(x, y, r, nil, style)
		x += runewidth.RuneWidth(r)
	}
	return x
}

// autoAlign returns a best-guess align for a Piece: the outermost indices are
// left- and right adjusted, while the rest are centered.
func autoAlign(index, total int) int {
	switch index {
	case 0:
		return AlignLeft
	case total - 1:
		return AlignRight
	default:
		return AlignCenter
	}
}

// getPiecesStartX calculates the start x-position for a given piece.
//
// Unused space is avoided by assigning extra space to the first pieces,
// if (xmax / pieces) leaves a remainder.
func getPiecesStartX(piece, pieces, xmax int) int {
	x := piece * (xmax / pieces)
	if piece <= (xmax % pieces) {
		return x + piece
	}
	return x + (xmax % pieces)
}

// alignX returns the draw start position.
func alignX(x, bufferWidth, textWidth, align int) int {
	switch align {
	case AlignLeft:
		return x
	case AlignCenter:
		return x + (bufferWidth / 2) - (textWidth / 2)
	case AlignRight:
		return x + bufferWidth - textWidth
	default:
		return x
	}
}

func pieceTextWidth(piece *topbar.PieceStatement) int {
	width := 0
	for _, fragment := range piece.Fragments {
		s, _ := fragment.Instance.Text()
		width += len(s)
	}
	return width
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
}
