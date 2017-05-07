package widgets

import (
	"fmt"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/message"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type MultibarWidget struct {
	// runes contain user input, while
	// text and errorText contains text set by the program.
	runes     []rune
	msg       message.Message
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
	MultibarModeNormal = iota
	MultibarModeVisual
	MultibarModeInput
	MultibarModeSearch
)

func NewMultibarWidget(events chan parser.KeyEvent) *MultibarWidget {
	m := &MultibarWidget{}
	m.runes = make([]rune, 0)
	m.events = events
	return m
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

func (m *MultibarWidget) SetMode(mode int) error {
	switch mode {
	case MultibarModeNormal:
	case MultibarModeVisual:
	case MultibarModeInput:
	case MultibarModeSearch:
	default:
		return fmt.Errorf("Mode not supported")
	}
	console.Log("Switching input mode from %d to %d", m.inputMode, mode)
	m.inputMode = mode
	m.setRunes([]rune{})
	PostEventModeSync(m, m.inputMode)
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
	case MultibarModeInput:
		s = ":" + m.RuneString()
		st = m.Style("commandText")
	case MultibarModeSearch:
		s = "/" + m.RuneString()
		st = m.Style("searchText")
	case MultibarModeVisual:
		s = "-- VISUAL --"
		st = m.Style("visualText")
	default:
		s = m.msg.Text
		st = m.textStyle
	}

	m.SetLeft(s, st)
}

func (m *MultibarWidget) RuneString() string {
	return string(m.runes)
}

func (m *MultibarWidget) RuneLen() int {
	return len(m.runes)
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
		m.SetMode(MultibarModeNormal)
	}
}

// handleTextInputEvent is called when an input event is received during any of the text input modes.
func (m *MultibarWidget) handleTextInputEvent(ev *tcell.EventKey) bool {
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
	return false
}

// handleNormalEvent is called when an input event is received during command mode.
func (m *MultibarWidget) handleNormalEvent(ev *tcell.EventKey) bool {
	ke := parser.KeyEvent{Key: ev.Key(), Rune: ev.Rune()}
	//console.Log("Input event in command mode: %s %s", ke.Key, string(ke.Rune))
	m.events <- ke
	return true
}

func (m *MultibarWidget) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch m.inputMode {
		case MultibarModeVisual:
			// Visual mode should have the same key bindings and behavior as normal mode.
			fallthrough
		case MultibarModeNormal:
			return m.handleNormalEvent(ev)
		case MultibarModeInput:
			return m.handleTextInputEvent(ev)
		case MultibarModeSearch:
			return m.handleTextInputEvent(ev)
		}
	}
	return false
}

func (w *MultibarWidget) Resize() {
}
