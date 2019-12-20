package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/style"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// Multibar receives keyboard events, displays status messages, and the position readout.
type Multibar struct {
	api  api.API
	view views.View
	views.WidgetWatchers
	style.Styled
}

var _ views.Widget = &Multibar{}

func NewMultibarWidget(a api.API) *Multibar {
	return &Multibar{
		api: a,
	}
}

func (w *Multibar) messageStyle(msg log.Message) tcell.Style {
	switch {
	case msg.Level == log.InfoLevel:
		return w.Style("statusbar")
	case msg.Level == log.ErrorLevel:
		return w.Style("errorText")
	default:
		return w.Style("default")
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
	hasVisualSelection := w.api.Songlist() != nil && w.api.Songlist().HasVisualSelection()
	sequenceText := w.api.Sequencer().String()
	multibarMode := w.api.Multibar().Mode()
	msg := log.Last(log.InfoLevel)

	switch {
	case multibarMode == multibar.ModeInput:
		return ":" + w.api.Multibar().String(), w.Style("commandText")
	case multibarMode == multibar.ModeSearch:
		return "/" + w.api.Multibar().String(), w.Style("searchText")
	case len(sequenceText) > 0:
		return sequenceText, w.Style("sequenceText")
	case hasVisualSelection:
		return "-- VISUAL --", w.Style("visualText")
	case msg != nil:
		return msg.Text, w.messageStyle(*msg)
	default:
		return "", w.Style("default")
	}
}

// Draw the statusbar part of the Multibar.
func (w *Multibar) Draw() {
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
	// text := w.api.SonglistWidget().PositionReadout()
	// FIXME
	text := ""
	x, _ := w.Size()
	x -= len(text)
	for _, r := range text {
		w.view.SetContent(x, 0, r, []rune{}, st)
		x++
	}

}
