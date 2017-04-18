package widgets

import (
	"fmt"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input/parser"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type MultibarWidget struct {
	// runes contain user input, while
	// text and errorText contains text set by the program.
	runes     []rune
	text      string
	textStyle tcell.Style
	events    chan parser.KeyEvent

	inputMode int
	styles    StyleMap

	views.TextBar
	widget
}

// Different input modes are handled in different ways. Check
// MultibarWidget.inputMode against these constants.
const (
	MultibarModeCommand = iota
	MultibarModeCommandInput
	MultibarModeSearch
)

func NewMultibarWidget(events chan parser.KeyEvent) *MultibarWidget {
	m := &MultibarWidget{}
	m.runes = make([]rune, 0)
	m.events = events
	return m
}

func (m *MultibarWidget) SetText(format string, a ...interface{}) {
	m.text = fmt.Sprintf(format, a...)
	m.textStyle = m.Style("statusbar")
	m.DrawStatusbar()
}

func (m *MultibarWidget) SetErrorText(format string, a ...interface{}) {
	m.text = fmt.Sprintf(format, a...)
	m.textStyle = m.Style("errorText")
	m.DrawStatusbar()
}

func (m *MultibarWidget) SetSequenceText(format string, a ...interface{}) {
	m.text = fmt.Sprintf(format, a...)
	m.textStyle = m.Style("sequenceText")
	m.DrawStatusbar()
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
	m.runes = r
	m.DrawStatusbar()
}

// Draw the statusbar part of the Multibar.
func (m *MultibarWidget) DrawStatusbar() {
	var st tcell.Style
	var s string

	switch m.inputMode {
	case MultibarModeCommandInput:
		s = ":" + m.RuneString()
		st = m.Style("commandText")
	case MultibarModeSearch:
		s = "/" + m.RuneString()
		st = m.Style("searchText")
	default:
		s = m.text
		st = m.textStyle
		m.text = ""
	}

	m.SetLeft(s, st)
}

func (m *MultibarWidget) RuneString() string {
	return string(m.runes)
}

func (m *MultibarWidget) handleTruncate() {
	m.setRunes(make([]rune, 0))
	PostEventInputChanged(m)
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
			m.handleTextRune(ev.Rune())
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
		ke := parser.KeyEvent{Key: ev.Key(), Rune: ev.Rune()}
		console.Log("Input event in command mode: %s %s", ke.Key, string(ke.Rune))
		m.events <- ke
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

func (w *MultibarWidget) Resize() {
}
