package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/style"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"time"
)

// Multibar receives keyboard events, displays status messages, and the position readout.
type Multibar struct {
	api              api.API
	view             views.View
	messageTimestamp time.Time
	views.WidgetWatchers
	style.Styled
}

var _ views.Widget = &Multibar{}

func NewMultibarWidget(a api.API) *Multibar {
	return &Multibar{
		api: a,
	}
}

func (w *Multibar) SetView(view views.View) {
	w.view = view
}

func (w *Multibar) Size() (int, int) {
	x, _ := w.view.Size()
	return x, 1
}

func (w *Multibar) Resize() {
}

func (w *Multibar) HandleEvent(ev tcell.Event) bool {
	return false
}

// Figure out what the multibar should render, and return it with the corresponding style
func (w *Multibar) textWithStyle() (string, tcell.Style) {
	hasVisualSelection := w.api.List() != nil && w.api.List().HasVisualSelection()
	sequenceText := w.api.Sequencer().String()
	multibarMode := w.api.Multibar().Mode()
	msg := log.Last(log.InfoLevel)

	switch {
	case multibarMode == multibar.ModeInput:
		w.messageTimestamp = time.Now()
		return ":" + w.api.Multibar().String(), w.Style("commandText")
	case multibarMode == multibar.ModeSearch:
		w.messageTimestamp = time.Now()
		return "/" + w.api.Multibar().String(), w.Style("searchText")
	case len(sequenceText) > 0:
		w.messageTimestamp = time.Now()
		return sequenceText, w.Style("sequenceText")
	case hasVisualSelection:
		w.messageTimestamp = time.Now()
		return "-- VISUAL --", w.Style("visualText")
	case msg != nil && w.messageTimestamp.UnixNano() < msg.Timestamp.UnixNano():
		return msg.Text, w.MessageStyle(*msg)
	default:
		return "", w.Style("default")
	}
}

// Draw the statusbar part of the Multibar.
func (w *Multibar) Draw() {
	w.SetStylesheet(w.api.Styles())
	w.view.Clear()
	w.drawLeft()
	w.drawRight()
}

// Draw the statusbar part of the Multibar.
func (w *Multibar) drawLeft() {
	text, st := w.textWithStyle()

	log.Debugf("multibar draw in style %x: %s", st, text)
	for x, r := range text {
		w.view.SetContent(x, 0, r, []rune{}, st)
		x++
	}
}

func (w *Multibar) drawRight() {
	st := w.Style("readout")
	// text := w.api.TableWidget().PositionReadout()
	// FIXME
	text := ""
	x, _ := w.Size()
	x -= len(text)
	for _, r := range text {
		w.view.SetContent(x, 0, r, []rune{}, st)
		x++
	}

}
