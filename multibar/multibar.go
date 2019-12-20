// package multibar provides a combined Vi-like statusbar and input field.
//
// Multibar has three modes:
//   * NORMAL	statusbar text is shown
//   * COMMAND	acts as a command input box
//   * SEARCH	acts as a search input box

package multibar

import (
	"github.com/ambientsound/pms/log"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/utils"

	"github.com/gdamore/tcell"
)

type TabCompleter interface {
	Scan() (string, error)
}

type TabCompleterFactory func(input string) TabCompleter

// Multibar implements a Vi-like combined statusbar and input box.
type Multibar struct {
	buffer      []rune
	commands    chan string
	cursor      int
	history     []*history
	mode        InputMode
	msg         message.Message
	searches    chan string
	tabComplete TabCompleter
	tcf         TabCompleterFactory
}

func New(tcf TabCompleterFactory) *Multibar {
	hist := make([]*history, 3)
	for i := range hist {
		hist[i] = NewHistory()
	}
	return &Multibar{
		history:  hist,
		buffer:   make([]rune, 0),
		commands: make(chan string, 1),
		searches: make(chan string, 1),
		tcf:      tcf,
	}
}

// Input is called on keyboard events.
func (m *Multibar) Input(event tcell.Event) bool {
	ev, ok := event.(*tcell.EventKey)
	if !ok {
		return false
	}

	if m.mode == ModeNormal {
		return false
	}

	log.Debugf("multibar keypress: name=%v key=%v modifiers=%v", ev.Name(), ev.Key(), ev.Modifiers())

	switch ev.Key() {

	// Alt keys has to be handled a bit differently than Ctrl keys.
	case tcell.KeyRune:
		modifiers := ev.Modifiers()
		if modifiers&tcell.ModAlt == 0 {
			// Pass the rune on to the text handling function if the alt modifier was not used.
			m.inputRune(ev.Rune())
		} else {
			switch ev.Rune() {
			case 'b':
				m.wordJump(-1)
			case 'f':
				m.wordJump(1)
			}
		}

	case tcell.KeyCtrlU:
		m.truncate()
	case tcell.KeyEnter:
		m.finish()
	case tcell.KeyTab:
		if m.Mode() == ModeInput {
			m.tab()
		}
	case tcell.KeyLeft, tcell.KeyCtrlB:
		m.moveCursor(-1)
	case tcell.KeyRight, tcell.KeyCtrlF:
		m.moveCursor(1)
	case tcell.KeyUp, tcell.KeyCtrlP:
		m.moveHistory(-1)
	case tcell.KeyDown, tcell.KeyCtrlN:
		m.moveHistory(1)
	case tcell.KeyCtrlG, tcell.KeyCtrlC:
		m.abort()
	case tcell.KeyCtrlA, tcell.KeyHome:
		m.moveCursor(-len(m.buffer))
	case tcell.KeyCtrlE, tcell.KeyEnd:
		m.moveCursor(len(m.buffer))
	case tcell.KeyBS, tcell.KeyDEL:
		m.backspace()
	case tcell.KeyCtrlW:
		m.deleteWord()

	default:
		console.Log("Unhandled text input event in Multibar: %v", ev.Key())
		return false
	}

	return true
}

// History returns the input history of the current input mode.
func (m *Multibar) History() *history {
	return m.history[m.mode]
}

// Clear the statusbar text
func (m *Multibar) Clear() {
	m.SetMessage(message.Message{
		Severity: message.Info,
		Text:     "",
	})
}

// Set an error in the statusbar
func (m *Multibar) Error(err error) {
	m.SetMessage(message.Message{
		Severity: message.Error,
		Text:     err.Error(),
	})
}

func (m *Multibar) Message() message.Message {
	return m.msg
}

func (m *Multibar) SetMessage(msg message.Message) {
	m.msg = msg
}

func (m *Multibar) SetMode(mode InputMode) {
	console.Log("Switching input mode from %s to %s", m.mode, mode)
	m.mode = mode
	m.setRunes(make([]rune, 0))
	m.History().Reset("")
}

func (m *Multibar) Mode() InputMode {
	return m.mode
}

func (m *Multibar) String() string {
	return string(m.buffer)
}

func (m *Multibar) Len() int {
	return len(m.buffer)
}

// Cursor returns the cursor position.
func (m *Multibar) Cursor() int {
	return m.cursor
}

// Commands returns a channel sending any commands entered in input mode.
func (m *Multibar) Commands() <-chan string {
	return m.commands
}

// Searches returns a channel sending any search terms.
func (m *Multibar) Searches() <-chan string {
	return m.searches
}

func (m *Multibar) setRunes(r []rune) {
	m.buffer = r
	m.validateCursor()
}

// validateCursor makes sure the cursor stays within boundaries.
func (m *Multibar) validateCursor() {
	if m.cursor > len(m.buffer) {
		m.cursor = len(m.buffer)
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *Multibar) truncate() {
	m.tabComplete = nil
	m.setRunes(make([]rune, 0))
	m.History().Reset(m.String())
}

// inputRune inserts a literal rune at the cursor position.
func (m *Multibar) inputRune(r rune) {
	m.tabComplete = nil
	runes := make([]rune, len(m.buffer)+1)
	copy(runes, m.buffer[:m.cursor])
	copy(runes[m.cursor+1:], m.buffer[m.cursor:])
	runes[m.cursor] = r
	m.setRunes(runes)

	m.cursor++
	m.History().Reset(m.String())
}

// backspace deletes a literal rune behind the cursor position.
func (m *Multibar) backspace() {

	m.tabComplete = nil

	// Backspace on an empty string returns to normal mode.
	if len(m.buffer) == 0 {
		m.abort()
		return
	}

	// Copy all runes except the deleted rune
	runes := deleteBackwards(m.buffer, m.cursor, 1)
	m.cursor--
	m.setRunes(runes)

	m.History().Reset(m.String())
}

// deleteWord deletes the previous word, along with all the backspace
// succeeding it.
func (m *Multibar) deleteWord() {

	m.tabComplete = nil

	// We don't use the lexer here because it is too smart when it comes to
	// quoted strings.
	cursor := m.cursor - 1

	// Scan backwards until a non-space character is found.
	for ; cursor >= 0; cursor-- {
		if !unicode.IsSpace(m.buffer[cursor]) {
			break
		}
	}

	// Scan backwards until a space character is found.
	for ; cursor >= 0; cursor-- {
		if unicode.IsSpace(m.buffer[cursor]) {
			cursor++
			break
		}
	}

	// Delete backwards.
	runes := deleteBackwards(m.buffer, m.cursor, m.cursor-cursor)
	m.cursor = cursor
	m.setRunes(runes)

	m.History().Reset(m.String())
}

func (m *Multibar) finish() {
	input := m.String()
	m.tabComplete = nil
	m.History().Add(input)

	mode := m.mode
	m.SetMode(ModeNormal)

	switch mode {
	case ModeInput:
		m.commands <- input
	case ModeSearch:
		m.searches <- input
	}
}

func (m *Multibar) abort() {
	m.setRunes(make([]rune, 0))
	m.finish()
}

func (m *Multibar) moveHistory(offset int) {
	m.tabComplete = nil
	s := m.History().Navigate(offset)
	m.setRunes([]rune(s))
	m.cursor = len(m.buffer)
}

func (m *Multibar) moveCursor(offset int) {
	m.tabComplete = nil
	m.cursor += offset
	m.validateCursor()
}

// wordJump moves the cursor forward to the start of the next word or
// backwards to the start of the previous word.
func (m *Multibar) wordJump(offset int) {
	m.tabComplete = nil
	m.cursor += nextWord(m.buffer, m.cursor, offset)
	m.validateCursor()
}

// tab invokes tab completion.
func (m *Multibar) tab() {

	// Ignore event if cursor is not at the end
	if m.cursor != len(m.buffer) {
		return
	}

	// Initialize tabcomplete
	if m.tabComplete == nil {
		m.tabComplete = m.tcf(m.String())
	}

	// Get next sentence, and abort on any errors.
	sentence, err := m.tabComplete.Scan()
	if err != nil {
		console.Log("Autocomplete: %s", err)
		return
	}

	// Replace current text.
	m.setRunes([]rune(sentence))
	m.cursor = len(m.buffer)
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
