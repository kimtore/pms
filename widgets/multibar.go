package widgets

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/constants"
	"github.com/ambientsound/pms/input/lexer"
	input_parser "github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/parser"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/utils"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

// history represents a text history that can be navigated through.
type history struct {
	items   []string
	current string
	index   int
}

// autocomplete represents a list of strings that autocompletes the current word.
type autocomplete struct {
	active   bool
	base     string
	index    int
	items    []string
	original string
}

// MultibarWidget receives keyboard events, displays status messages, and the position readout.
type MultibarWidget struct {
	api          api.API
	autocomplete autocomplete
	cursor       int
	events       chan input_parser.KeyEvent
	inputMode    int
	msg          message.Message
	runes        []rune
	textStyle    tcell.Style

	// Three histories, one for each input mode
	history [3]history

	views.TextBar
	style.Styled
}

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

func NewMultibarWidget(a api.API, events chan input_parser.KeyEvent) *MultibarWidget {
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
	case constants.MultibarModeNormal:
	case constants.MultibarModeInput:
	case constants.MultibarModeSearch:
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
	m.autocomplete.active = false
	m.setRunes(make([]rune, 0))
	m.History().Reset(m.RuneString())
	PostEventInputChanged(m)
}

// handleTextRune inserts a literal rune at the cursor position.
func (m *MultibarWidget) handleTextRune(r rune) {
	m.autocomplete.active = false
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

	m.autocomplete.active = false

	// Backspace on an empty string returns to normal mode.
	if len(m.runes) == 0 {
		m.SetMode(constants.MultibarModeNormal)
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
	m.autocomplete.active = false
	m.History().Add(m.RuneString())
	PostEventInputFinished(m)
}

func (m *MultibarWidget) handleAbort() {
	m.autocomplete.active = false
	m.History().Add(m.RuneString())
	m.History().Reset("")
	m.setRunes(make([]rune, 0))
	m.SetMode(constants.MultibarModeNormal)
}

func (m *MultibarWidget) handleHistory(offset int) {
	m.autocomplete.active = false
	s := m.History().Navigate(offset)
	m.setRunes([]rune(s))
	m.cursor = len(m.runes)
	PostEventInputChanged(m)
}

func (m *MultibarWidget) handleCursor(offset int) {
	m.autocomplete.active = false
	m.cursor += offset
	m.validateCursor()
	PostEventInputChanged(m) // FIXME: this triggers a search query; disable that
}

// nextWord returns the distance to the next word in a rune slice.
func nextWord(runes []rune, cursor, offset int) int {
	var s string

	switch {
	// Move backwards
	case offset < 0:
		rev := utils.ReverseRunes(runes)
		revIndex := len(runes) - cursor
		runes := rev[revIndex:]
		s = string(runes)

	// Move forwards
	case offset > 0:
		runes := runes[cursor:]
		s = string(runes)

	default:
		return 0
	}

	reader := strings.NewReader(s)
	scanner := lexer.NewScanner(reader)

	// Strip any whitespace, and count the total length of the whitespace and
	// the next token.
	tok, lit := scanner.Scan()
	skip := utf8.RuneCountInString(lit)
	if tok == lexer.TokenWhitespace {
		_, lit = scanner.Scan()
		skip += utf8.RuneCountInString(lit)
	}

	return offset * skip
}

// handleCursorWord moves the cursor forward to the start of the next word or
// backwards to the start of the previous word.
func (m *MultibarWidget) handleCursorWord(offset int) {
	m.autocomplete.active = false
	m.cursor += nextWord(m.runes, m.cursor, offset)
	m.validateCursor()
	PostEventInputChanged(m) // FIXME: this triggers a search query; disable that
}

// handleTab cycles through autocomplete entries.
func (m *MultibarWidget) handleTab() {

	// Update text if autocomplete has been initialized
	if m.autocomplete.active {
		if len(m.autocomplete.items) == 0 {
			return
		}
		if m.autocomplete.index >= len(m.autocomplete.items) {
			m.autocomplete.index = 0
		}
		s := m.autocomplete.base + m.autocomplete.items[m.autocomplete.index]
		m.autocomplete.index++
		m.setRunes([]rune(s))
		m.cursor = len(m.runes)
		PostEventInputChanged(m) // FIXME: this triggers a search query; disable that
		return
	}

	// Set up the input token stream for the parser
	s := string(m.runes)
	reader := strings.NewReader(s)
	scanner := lexer.NewScanner(reader)
	parser := parser.New(scanner)

	// Find the verb
	tok, verb := parser.ScanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		console.Log("Tab completing verb '%s', but this is not an identifier", verb)
		return
	}

	// Instantiate the Command registered with this verb
	cmd := commands.New(verb, m.api)
	if cmd == nil {
		console.Log("Tab completing verb '%s' yielded zero results", verb)
		return
	}

	// Parse the remaining text
	cmd.Parse(scanner)

	// Concatenate scanned tokens, except the last one
	// FIXME: add support for command completion
	tokens := cmd.Scanned()
	if len(tokens) < 2 {
		return
	}

	//console.Log("Scanned tokens: %+v", tokens)

	lastToken := lexer.TokenEnd
	stringTokens := make([]string, 1)
	stringTokens[0] = verb
	for i := range tokens {
		if tokens[i].Tok == lexer.TokenEnd {
			break
		}
		lastToken = tokens[i].Tok
		stringTokens = append(stringTokens, tokens[i].Lit)
	}

	//console.Log("StringTokens: %+v", stringTokens)

	// Initialize autocomplete
	lt := len(stringTokens) - 1
	m.autocomplete.active = true
	m.autocomplete.index = len(m.autocomplete.items)
	m.autocomplete.items = cmd.TabComplete()
	if lastToken == lexer.TokenWhitespace {
		m.autocomplete.original = ""
		m.autocomplete.base = strings.Join(stringTokens, "")
	} else {
		m.autocomplete.original = stringTokens[lt]
		m.autocomplete.base = strings.Join(stringTokens[:lt], "")
	}

	// Recurse to update
	m.handleTab()
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

	// Alt keys has to be handled a bit differently than Ctrl keys.
	case tcell.KeyRune:
		modifiers := ev.Modifiers()
		if modifiers&tcell.ModAlt == 0 {
			// Pass the rune on to the text handling function if the alt modifier was not used.
			m.handleTextRune(ev.Rune())
		} else {
			switch ev.Rune() {
			case 'b':
				m.handleCursorWord(-1)
			case 'f':
				m.handleCursorWord(1)
			}
		}

	case tcell.KeyCtrlU:
		m.handleTruncate()
	case tcell.KeyEnter:
		m.handleFinished()
	case tcell.KeyTab:
		if m.Mode() == constants.MultibarModeInput {
			m.handleTab()
		}
	case tcell.KeyLeft, tcell.KeyCtrlB:
		m.handleCursor(-1)
	case tcell.KeyRight, tcell.KeyCtrlF:
		m.handleCursor(1)
	case tcell.KeyUp, tcell.KeyCtrlP:
		m.handleHistory(-1)
	case tcell.KeyDown, tcell.KeyCtrlN:
		m.handleHistory(1)
	case tcell.KeyCtrlG, tcell.KeyCtrlC:
		m.handleAbort()
	case tcell.KeyCtrlA, tcell.KeyHome:
		m.handleCursor(-len(m.runes))
	case tcell.KeyCtrlE, tcell.KeyEnd:
		m.handleCursor(len(m.runes))
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
	ke := input_parser.KeyEvent{Key: ev.Key(), Rune: ev.Rune()}
	//console.Log("Input event in command mode: %s %s", ke.Key, string(ke.Rune))
	m.events <- ke
	return true
}

func (m *MultibarWidget) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch m.inputMode {
		case constants.MultibarModeNormal:
			return m.handleNormalEvent(ev)
		case constants.MultibarModeInput:
			return m.handleTextInputEvent(ev)
		case constants.MultibarModeSearch:
			return m.handleTextInputEvent(ev)
		}
	}
	return false
}

func (w *MultibarWidget) Resize() {
}
