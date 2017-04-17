package widgets

import (
	"github.com/gdamore/tcell"
)

type widget struct {
	styles StyleMap
}

func (w *widget) Style(s string) tcell.Style {
	return w.styles[s]
}

func (w *widget) SetStyleMap(m StyleMap) {
	w.styles = m
}
