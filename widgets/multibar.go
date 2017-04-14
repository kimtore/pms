package widgets

import (
	"fmt"

	"github.com/ambientsound/pms/console"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type MultibarWidget struct {
	views.TextBar

	runes       []rune
	defaultText string
	inputMode   int
}

const (
	MultibarModeCommand = iota
	MultibarModeCommandInput
	MultibarModeSearch
)

func NewMultibarWidget() *MultibarWidget {
	m := &MultibarWidget{}
	m.runes = make([]rune, 0)
	return m
}

func (m *MultibarWidget) SetDefaultText(s string) {
	m.defaultText = s
	m.setRunes(m.runes)
}

func (m *MultibarWidget) SetMode(mode int) error {
	switch mode {
	case MultibarModeCommand:
	case MultibarModeCommandInput:
	case MultibarModeSearch:
	default:
		return fmt.Errorf("Mode not supported")
	}
	console.Log("Switching input mode from %d to %d", m.inputMode, mode)
	m.inputMode = mode
	m.setRunes([]rune{})
	PostEventInputChanged(m)
	return nil
}

func (m *MultibarWidget) Mode() int {
	return m.inputMode
}

func (m *MultibarWidget) setRunes(r []rune) {
	var s string
	m.runes = r
	s = m.RuneString()

	// Visual feedback
	switch m.inputMode {
	case MultibarModeCommandInput:
		s = ":" + s
	case MultibarModeSearch:
		s = "/" + s
	default:
		s = m.defaultText
	}

	m.SetLeft(s, tcell.StyleDefault)
	//console.Log("Multibar> %s", s)
}

func (m *MultibarWidget) RuneString() string {
	return string(m.runes)
}

func (m *MultibarWidget) handleTruncate() {
	m.setRunes(make([]rune, 0))
	PostEventInputChanged(m)
}

func (m *MultibarWidget) handleRune(r rune) {
	switch m.inputMode {
	case MultibarModeCommand:
		switch r {
		case '/':
			m.SetMode(MultibarModeSearch)
		case ':':
			m.SetMode(MultibarModeCommandInput)
		default:
			console.Log("Unhandled input rune: %s", string(r))
		}
	case MultibarModeCommandInput:
		fallthrough
	case MultibarModeSearch:
		m.handleTextRune(r)
	}
}

func (m *MultibarWidget) handleTextRune(r rune) {
	m.setRunes(append(m.runes, r))
	PostEventInputChanged(m)
}

func (m *MultibarWidget) handleBackspace() {
	if len(m.runes) > 0 {
		m.setRunes(m.runes[:len(m.runes)-1])
		PostEventInputChanged(m)
	} else {
		m.SetMode(MultibarModeCommand)
	}
}

// handleTextInputEvent is called when an input event is received during any of the text input modes.
func (m *MultibarWidget) handleTextInputEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			m.handleRune(ev.Rune())
			return true
		case tcell.KeyCtrlU:
			m.handleTruncate()
			return true
		case tcell.KeyEnter:
			PostEventInputFinished(m)
			return true
		case tcell.KeyBS:
			fallthrough
		case tcell.KeyDEL:
			m.handleBackspace()
			return true
		}
		console.Log("Unhandled text input event in Multibar: %s", ev.Key())
	}
	return false
}

// handleCommandEvent is called when an input event is received during command mode.
func (m *MultibarWidget) handleCommandEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			switch ev.Rune() {
			case '/':
				m.SetMode(MultibarModeSearch)
				return true
			case ':':
				m.SetMode(MultibarModeCommandInput)
				return true
			}
		}
		console.Log("Unhandled input event in command mode: %s %s", ev.Key(), string(ev.Rune()))
	}
	return false
}

func (m *MultibarWidget) HandleEvent(ev tcell.Event) bool {
	switch m.inputMode {
	case MultibarModeCommand:
		return m.handleCommandEvent(ev)
	case MultibarModeCommandInput:
		return m.handleTextInputEvent(ev)
	case MultibarModeSearch:
		return m.handleTextInputEvent(ev)
	}
	return false
}
