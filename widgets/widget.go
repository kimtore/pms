package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type widget struct {
	styles StyleMap
	views.WidgetWatchers
}

func (w *widget) SetStyleMap(m StyleMap) {
	w.styles = m
}

func (w *widget) Resize() {
}

func (w *widget) HandleEvent(ev tcell.Event) bool {
	return false
}
