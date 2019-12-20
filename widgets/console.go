package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/log"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// ConsoleWidget is a tcell widget which draws the program log.
type ConsoleWidget struct {
	api      api.API
	view     views.View
	viewport views.ViewPort
	views.WidgetWatchers
}

var _ views.Widget = &ConsoleWidget{}

func NewConsoleWidget() *ConsoleWidget {
	return &ConsoleWidget{}
}

func (w *ConsoleWidget) SetView(view views.View) {
	w.view = view
	w.viewport.SetView(view)
	log.Debugf("console widget: set view %#v", view)
}

func (w *ConsoleWidget) Size() (int, int) {
	x, y := w.view.Size()
	log.Debugf("console widget: report size %d x %d", x, y)
	return w.view.Size()
}

func (w *ConsoleWidget) Draw() {
	log.Debugf("console widget: draw")

	list := log.Messages(log.InfoLevel)
	entries := len(list)
	_, ymax := w.Size()
	if entries > ymax {
		list = list[entries-ymax:]
	}

	w.viewport.Clear()
	st := tcell.StyleDefault

	for y, msg := range list {
		x := 0
		ts := msg.Timestamp.Format(time.RFC822)
		x = w.drawString(x, y, ts, st)
		x = w.drawString(x+1, y, msg.Level.String(), st)
		x = w.drawString(x+1, y, msg.Text, st)
	}
}

func (w *ConsoleWidget) Resize() {
	log.Debugf("console widget: resize")
	w.viewport.Resize(0, 0, -1, -1)
}

func (w *ConsoleWidget) HandleEvent(ev tcell.Event) bool {
	log.Debugf("console event: %#v", ev)
	return false
}

func (w *ConsoleWidget) drawString(x, y int, s string, style tcell.Style) int {
	for _, r := range s {
		w.view.SetContent(x, y, r, nil, style)
		x++
	}
	return x
}
