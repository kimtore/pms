package widgets

import (
	"github.com/ambientsound/pms/console"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type MultibarWidget struct {
	views.TextBar

	runes       []rune
	defaultText string
}

func NewMultibarWidget() *MultibarWidget {
	m := &MultibarWidget{}
	m.runes = make([]rune, 0)
	return m
}

func (m *MultibarWidget) SetDefaultText(s string) {
	m.defaultText = s
	m.setRunes(m.runes)
}

func (m *MultibarWidget) setRunes(r []rune) {
	var s string
	m.runes = r
	if len(m.runes) > 0 {
		s = m.GetRuneString()
	} else {
		s = m.defaultText
	}
	m.SetLeft(s, tcell.StyleDefault)
	console.Log("Multibar> %s", s)
}

func (m *MultibarWidget) GetRuneString() string {
	return string(m.runes)
}

func (m *MultibarWidget) handleTruncate() {
	m.setRunes(make([]rune, 0))
	PostEventInputChanged(m)
}

func (m *MultibarWidget) handleInputRune(r rune) {
	m.setRunes(append(m.runes, r))
	PostEventInputChanged(m)
}

func (m *MultibarWidget) handleBackspace() {
	if len(m.runes) > 0 {
		m.setRunes(m.runes[:len(m.runes)-1])
		PostEventInputChanged(m)
	}
}

func (m *MultibarWidget) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {

		case tcell.KeyRune:
			m.handleInputRune(ev.Rune())
			return true

		case tcell.KeyCtrlU:
			fallthrough
		case tcell.KeyEnter:
			m.handleTruncate()
			return true

		case tcell.KeyBS:
			fallthrough
		case tcell.KeyDEL:
			m.handleBackspace()
			return true

		}
		console.Log("%s", ev.Key())
	}
	return false
}
