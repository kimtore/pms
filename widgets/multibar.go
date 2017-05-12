package widgets

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/style"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// history represents a text history that can be navigated through.
type history struct {
	items   []string
	current string
	index   int
}

// MultibarWidget receives keyboard events, displays status messages, and the position readout.
type MultibarWidget struct {
	api       api.API
	cursor    int
	events    chan parser.KeyEvent
	inputMode int
	msg       message.Message
	runes     []rune
	textStyle tcell.Style

	// Three histories, one for each input mode
	history [3]history

	views.TextBar
	style.Styled
}

// Different input modes are handled in different ways. Check
// MultibarWidget.inputMode against these constants.
const (
	MultibarModeNormal = iota
	MultibarModeInput
	MultibarModeSearch
)

// Add adds to the input history.
func (h *history) Add(s string) {
	if len(s) > 0 {
		hl := len(h.items)
		if hl == 0 || h.items[hl-1] != s {
			h.items = append(h.items, s)
		}
	}
	h.Reset(s)
}

// Reset resets the cursor offset to the last position.
func (h *history) Reset(s string) {
	h.index = len(h.items)
	h.current = s
}

// Current returns the current history item.
func (h *history) Current() string {
	if len(h.items) == 0 || h.index >= len(h.items) {
		console.Log("Want index %d, returning current string '%s'", h.index, h.current)
		h.index = len(h.items)
		return h.current
	}
	h.validateIndex()
	console.Log("History returning index %d", h.index)
	return h.items[h.index]
}

// Navigate navigates the history and returns that history item.
func (h *history) Navigate(offset int) string {
	h.index += offset
	return h.Current()
}

// validateIndex ensures that the item index stays within the valid range.
func (h *history) validateIndex() {
	if h.index >= len(h.items) {
		h.index = len(h.items) - 1
	}
	if h.index < 0 {
		h.index = 0
	}
}

func NewMultibarWidget(a api.API, events chan parser.KeyEvent) *MultibarWidget {
	return &MultibarWidget{
		api:    a,
		runes:  make([]rune, 0),
		events: events,
		history: [3]history{
			{items: make([]string, 0)},
			{items: make([]string, 0)},
			{items: make([]string, 0)},
		},
	}
}

// History returns the input history of the current input mode.
func (m *MultibarWidget) History() *history {
	return &m.history[m.inputMode]
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
	case MultibarModeInput:
	case MultibarModeSearch:
	default:
		return fmt.Errorf("Mode not supported")
	}
	console.Log("Switching input mode from %d to %d", m.inputMode, mode)
	m.inputMode = mode
	m.setRunes(make([]rune, 0))
	m.History().Reset("")
	PostEventInputChanged(m)
	return nil
}

func (m *MultibarWidget) Mode() int {
	return m.inputMode
}

func (m *MultibarWidget) setRunes(r []rune) {
	m.runes = r
	m.validateCursor()
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
	default:
		if len(m.msg.Text) == 0 && m.api.Songlist().HasVisualSelection() {
			s = "-- VISUAL --"
			st = m.Style("visualText")
		} else {
			s = m.msg.Text
			st = m.textStyle
		}
	}

	m.SetLeft(s, st)
}

func (m *MultibarWidget) RuneString() string {
	return string(m.runes)
}

func (m *MultibarWidget) RuneLen() int {
	return len(m.runes)
}

func (m *MultibarWidget) Cursor() int {
	return m.cursor
}

func (m *MultibarWidget) handleTruncate() {
	m.setRunes(make([]rune, 0))
	m.History().Reset(m.RuneString())
	PostEventInputChanged(m)
}

// handleTextRune inserts a literal rune at the cursor position.
func (m *MultibarWidget) handleTextRune(r rune) {
	runes := make([]rune, len(m.runes)+1)
	copy(runes, m.runes[:m.cursor])
	copy(runes[m.cursor+1:], m.runes[m.cursor:])
	runes[m.cursor] = r
	m.setRunes(runes)

	m.cursor++
	m.History().Reset(m.RuneString())
	PostEventInputChanged(m)
}

// deleteBackwards returns a new rune slice with a part cut out. If the deleted
// part is bigger than the string contains, deleteBackwards removes as much as
// possible.
func deleteBackwards(src []rune, cursor int, length int) []rune {
	if cursor < length {
		length = cursor
	}
	runes := make([]rune, len(src)-length)
	index := copy(runes, src[:cursor-length])
	copy(runes[index:], src[cursor:])
	return runes
}

// handleBackspace deletes a literal rune behind the cursor position.
func (m *MultibarWidget) handleBackspace() {

	// Backspace on an empty string returns to normal mode.
	if len(m.runes) == 0 {
		m.SetMode(MultibarModeNormal)
		return
	}

	// Copy all runes except the deleted rune
	runes := deleteBackwards(m.runes, m.cursor, 1)
	m.cursor--
	m.setRunes(runes)

	m.History().Reset(m.RuneString())
	PostEventInputChanged(m)
}

func (m *MultibarWidget) handleFinished() {
	m.History().Add(m.RuneString())
	PostEventInputFinished(m)
}

func (m *MultibarWidget) handleAbort() {
	m.History().Add(m.RuneString())
	m.History().Reset("")
	m.setRunes(make([]rune, 0))
	m.SetMode(MultibarModeNormal)
}

func (m *MultibarWidget) handleHistory(offset int) {
	s := m.History().Navigate(offset)
	m.setRunes([]rune(s))
	m.cursor = len(m.runes)
	PostEventInputChanged(m)
}

func (m *MultibarWidget) handleCursor(offset int) {
	m.cursor += offset
	m.validateCursor()
	PostEventInputChanged(m) // FIXME: this triggers a search query; disable that
}

// validateCursor makes sure the cursor stays within boundaries.
func (m *MultibarWidget) validateCursor() {
	if m.cursor > len(m.runes) {
		m.cursor = len(m.runes)
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// handleTextInputEvent is called when an input event is received during any of the text input modes.
func (m *MultibarWidget) handleTextInputEvent(ev *tcell.EventKey) bool {
	switch ev.Key() {

	case tcell.KeyRune:
		m.handleTextRune(ev.Rune())
	case tcell.KeyCtrlU:
		m.handleTruncate()
	case tcell.KeyEnter:
		m.handleFinished()
	case tcell.KeyLeft:
		m.handleCursor(-1)
	case tcell.KeyRight:
		m.handleCursor(1)
	case tcell.KeyUp:
		m.handleHistory(-1)
	case tcell.KeyDown:
		m.handleHistory(1)
	case tcell.KeyCtrlG, tcell.KeyCtrlC:
		m.handleAbort()
	case tcell.KeyBS, tcell.KeyDEL:
		m.handleBackspace()

	default:
		console.Log("Unhandled text input event in Multibar: %s", ev.Key())
		return false
	}

	return true
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
