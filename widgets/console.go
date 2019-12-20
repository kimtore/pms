package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/utils"

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
	w.viewport.Fill('-', tcell.StyleDefault)
	w.drawNext(0, 0, 10, 10, []rune("foobar"), tcell.StyleDefault)
}

func (w *ConsoleWidget) Resize() {
	log.Debugf("console widget: resize")
	w.viewport.Resize(0, 0, -1, -1)
}

func (w *ConsoleWidget) HandleEvent(ev tcell.Event) bool {
	log.Debugf("console event: %#v", ev)
	return false
}

func (w *ConsoleWidget) drawNext(x, y, strmin, strmax int, runes []rune, style tcell.Style) int {
	strmin = utils.Min(len(runes), strmin)
	n := 0
	for n < strmin {
		w.view.SetContent(x, y, runes[n], nil, style)
		n++
		x++
	}
	for n < strmax {
		w.view.SetContent(x, y, ' ', nil, style)
		n++
		x++
	}
	return x
}
