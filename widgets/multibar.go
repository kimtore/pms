package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/style"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// MultibarWidget receives keyboard events, displays status messages, and the position readout.
type MultibarWidget struct {
	api      api.API
	multibar *multibar.Multibar

	views.TextBar
	style.Styled
}

var _ views.Widget = &MultibarWidget{}

func NewMultibarWidget(a api.API, m *multibar.Multibar) *MultibarWidget {
	return &MultibarWidget{
		api:      a,
		multibar: m,
	}
}

func (m *MultibarWidget) messageStyle(msg message.Message) tcell.Style {
	switch {
	case msg.Type == message.SequenceText:
		return m.Style("sequenceText")
	case msg.Severity == message.Info:
		return m.Style("statusbar")
	case msg.Severity == message.Error:
		return m.Style("errorText")
	default:
		return m.Style("default")
	}
}

// Draw the statusbar part of the Multibar.
func (m *MultibarWidget) Render() {
	var st tcell.Style
	var s string

	switch m.multibar.Mode() {
	case multibar.ModeInput:
		s = ":" + m.multibar.String()
		st = m.Style("commandText")
	case multibar.ModeSearch:
		s = "/" + m.multibar.String()
		st = m.Style("searchText")
	default:
		msg := m.multibar.Message()
		if len(msg.Text) == 0 && m.api.Songlist() != nil && m.api.Songlist().HasVisualSelection() {
			s = "-- VISUAL --"
			st = m.Style("visualText")
		} else {
			s = msg.Text
			st = m.messageStyle(msg)
		}
	}

	m.SetLeft(s, st)
}
