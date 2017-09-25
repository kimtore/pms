package widgets

import (
	"github.com/ambientsound/pms/term"
)

type canvas struct {
	c     term.Canvas
	dirty bool
}

// SetCanvas provides a new drawing area for the widget.
func (w *canvas) SetCanvas(c term.Canvas) {
	w.c = c
	w.SetDirty(true)
}

func (w *canvas) SetDirty(dirty bool) {
	w.dirty = dirty
}
