// package bufin provides buffered text input with cursor, readline, tab completion, and history.
package bufin

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/history"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/tabcomplete"
	"github.com/ambientsound/pms/term"
	"github.com/ambientsound/pms/utils"
	termbox "github.com/nsf/termbox-go"
)

// State signifies what state the buffer is in.
type State int

// State constants.
const (
	StateBuffer State = iota // buffer is updated.
	StateCursor              // cursor is moved.
	StateReturn              // user presses <Return> during buffered mode.
	StateCancel              // user cancels buffered mode.
)

// Buffer buffers text input in a readline-like editing environment, and
// facilitates tab completion, history, and keeps a cursor.
type Buffer struct {
	api         api.API
	cursor      int              // cursor position
	history     *history.History // backlog of text input
	runes       []rune           // text input
	state       State
	tabComplete *tabcomplete.TabComplete // tab completion list
}

// New returns Buffer.
func New(api api.API) *Buffer {
	return &Buffer{
		api:     api,
		history: history.New(),
		runes:   make([]rune, 0),
	}
}

func (b *Buffer) String() string {
	return string(b.runes)
}

func (b *Buffer) Len() int {
	return len(b.runes)
}

// Cursor returns the cursor position.
func (b *Buffer) Cursor() int {
	return b.cursor
}

func (b *Buffer) State() State {
	return b.state
}

func (b *Buffer) setState(s State) {
	b.state = s
}

func (b *Buffer) setRunes(r []rune) {
	b.runes = r
	b.validateCursor()
	b.setState(StateBuffer)
}

// validateCursor makes sure the cursor stays within boundaries.
func (b *Buffer) validateCursor() {
	if b.cursor > b.Len() {
		b.cursor = b.Len()
	}
	if b.cursor < 0 {
		b.cursor = 0
	}
}

func (b *Buffer) handleTruncate() {
	b.tabComplete = nil
	b.setRunes(make([]rune, 0))
	b.history.Reset(b.String())
}

// handleTextRune inserts a literal rune at the cursor position.
func (b *Buffer) handleTextRune(r rune) {
	b.tabComplete = nil
	runes := make([]rune, len(b.runes)+1)
	copy(runes, b.runes[:b.cursor])
	copy(runes[b.cursor+1:], b.runes[b.cursor:])
	runes[b.cursor] = r
	b.setRunes(runes)

	b.cursor++
	b.history.Reset(b.String())
}

// handleBackspace deletes a literal rune behind the cursor position.
func (b *Buffer) handleBackspace() {

	b.tabComplete = nil

	// Backspace on an empty string returns to normal mode.
	if len(b.runes) == 0 {
		b.handleAbort()
		return
	}

	// Copy all runes except the deleted rune
	runes := deleteBackwards(b.runes, b.cursor, 1)
	b.cursor--
	b.setRunes(runes)

	b.history.Reset(b.String())
}

// handleDeleteWord deletes the previous word, along with all the backspace
// succeeding it.
func (b *Buffer) handleDeleteWord() {

	b.tabComplete = nil

	// We don't use the lexer here because it is too smart when it comes to
	// quoted strings.
	cursor := b.cursor - 1

	// Scan backwards until a non-space character is found.
	for ; cursor >= 0; cursor-- {
		if !unicode.IsSpace(b.runes[cursor]) {
			break
		}
	}

	// Scan backwards until a space character is found.
	for ; cursor >= 0; cursor-- {
		if unicode.IsSpace(b.runes[cursor]) {
			cursor++
			break
		}
	}

	// Delete backwards.
	runes := deleteBackwards(b.runes, b.cursor, b.cursor-cursor)
	b.cursor = cursor
	b.setRunes(runes)

	b.history.Reset(b.String())
}

func (b *Buffer) handleFinished() {
	b.tabComplete = nil
	b.history.Add(b.String())
}

func (b *Buffer) handleAbort() {
	b.setRunes(make([]rune, 0))
	b.handleFinished()
	b.setState(StateCancel)
}

func (b *Buffer) handleComplete() {
	b.handleFinished()
	b.setState(StateReturn)
}

func (b *Buffer) handleHistory(offset int) {
	b.tabComplete = nil
	s := b.history.Navigate(offset)
	b.setRunes([]rune(s))
	b.cursor = len(b.runes)
}

func (b *Buffer) handleCursor(offset int) {
	b.tabComplete = nil
	b.cursor += offset
	b.setState(StateCursor)
	b.validateCursor()
}

// handleCursorWord moves the cursor forward to the start of the next word or
// backwards to the start of the previous word.
func (b *Buffer) handleCursorWord(offset int) {
	b.tabComplete = nil
	b.cursor += nextWord(b.runes, b.cursor, offset)
	b.setState(StateCursor)
	b.validateCursor()
}

// handleTab invokes tab completion.
func (b *Buffer) handleTab() {

	// Ignore event if cursor is not at the end
	if b.cursor != len(b.runes) {
		return
	}

	// Initialize tabcomplete
	if b.tabComplete == nil {
		b.tabComplete = tabcomplete.New(b.String(), b.api)
	}

	// Get next sentence, and abort on any errors.
	sentence, err := b.tabComplete.Scan()
	if err != nil {
		console.Log("Autocomplete: %s", err)
		return
	}

	// Replace current text.
	b.setRunes([]rune(sentence))
	b.cursor = len(b.runes)
}

// handleTextInputEvent is called when an input event is received during any of the text input modes.
func (b *Buffer) handleTextInputEvent(ev term.KeyPress) bool {
	switch ev.Key {

	// Alt keys has to be handled a bit differently than Ctrl keys.
	case 0:
		if ev.Mod&term.ModAlt == 0 {
			// Pass the rune on to the text handling function if the alt modifier was not used.
			b.handleTextRune(ev.Ch)
		} else {
			switch ev.Ch {
			case 'b':
				b.handleCursorWord(-1)
			case 'f':
				b.handleCursorWord(1)
			}
		}

	case termbox.KeyCtrlU:
		b.handleTruncate()
	case termbox.KeyEnter:
		b.handleComplete()
	case termbox.KeyTab:
		b.handleTab()
	case termbox.KeyArrowLeft, termbox.KeyCtrlB:
		b.handleCursor(-1)
	case termbox.KeyArrowRight, termbox.KeyCtrlF:
		b.handleCursor(1)
	case termbox.KeyArrowUp, termbox.KeyCtrlP:
		b.handleHistory(-1)
	case termbox.KeyArrowDown, termbox.KeyCtrlN:
		b.handleHistory(1)
	case termbox.KeyCtrlG, termbox.KeyCtrlC:
		b.handleAbort()
	case termbox.KeyCtrlA, termbox.KeyHome:
		b.handleCursor(-len(b.runes))
	case termbox.KeyCtrlE, termbox.KeyEnd:
		b.handleCursor(len(b.runes))
	case termbox.KeyBackspace, termbox.KeyDelete:
		b.handleBackspace()
	case termbox.KeyCtrlW:
		b.handleDeleteWord()

	default:
		console.Log("Unhandled text input event: %+v", ev)
		return false
	}

	return true
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
