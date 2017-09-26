package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/term"

	"github.com/gdamore/tcell/views"
)

// MultibarWidget receives keyboard events, displays status messages, and the position readout.
type MultibarWidget struct {
	api api.API
	msg message.Message

	textStyle term.Style
	views.TextBar
	style.Styled
}

func NewMultibarWidget(a api.API) *MultibarWidget {
	return &MultibarWidget{
		api: a,
	}
}

func (m *MultibarWidget) SetMessage(msg message.Message) {
	switch {
	case msg.Type == message.SequenceText:
		m.textStyle = m.Style("sequenceText")
	case msg.Severity == message.Info:
		m.textStyle = m.Style("statusbar")
	case msg.Severity == message.Error:
		m.textStyle = m.Style("errorText")
	default:
		return
	}
	m.msg = msg
	m.DrawStatusbar()
}

// Draw the statusbar part of the Multibar.
func (m *MultibarWidget) DrawStatusbar() {
	/*
		var st term.Style
		var s string

		switch m.inputMode {
		case constants.MultibarModeInput:
			s = ":" + m.RuneString()
			st = m.Style("commandText")
		case constants.MultibarModeSearch:
			s = "/" + m.RuneString()
			st = m.Style("searchText")
		default:
			if len(m.msg.Text) == 0 && m.api.Songlist().HasVisualSelection() {
				s = "-- VISUAL --"
				st = m.Style("visualText")
			} else {
				s = m.msg.Text
				st = m.textStyle
			}
		}

		_, _ = s, st

		// FIXME: m.SetLeft(s, st)
	*/
}

func (w *MultibarWidget) Resize() {
}

func (w *MultibarWidget) SetMode(int) error {
	return nil
}

func (w *MultibarWidget) Mode() int {
	return 0
}
