package keysequence

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/parser"
	"github.com/gdamore/tcell"
)

// modifiers maps modifier names to their integer representation.
var modifiers = map[string]tcell.ModMask{
	"s":     tcell.ModShift,
	"c":     tcell.ModCtrl,
	"a":     tcell.ModAlt,
	"m":     tcell.ModMeta,
	"shift": tcell.ModShift,
	"ctrl":  tcell.ModCtrl,
	"alt":   tcell.ModAlt,
	"meta":  tcell.ModMeta,
}

// Parser is a key sequence parser. Key sequences consists of one or more
// keyboard input events. Special keys are defined by using angle brackets and
// a hyphen, such as <C-a> for Ctrl+A, or <C-S-M-x> for Ctrl+Shift+Meta+X.
type Parser struct {
	parser.Parser
}

// NewParser returns Parser.
func NewParser(r *lexer.Scanner) *Parser {
	return &Parser{
		parser.Parser{S: r},
	}
}

// ParseKeySequence parses the next key sequence, which is a combination of literal keys.
func (p *Parser) ParseKeySequence() (KeySequence, error) {

	keyseq := make(KeySequence, 0)

	tok, lit := p.ScanIgnoreWhitespace()

Parse:
	for {
		switch tok {

		// A left angle bracket signifies a special key, such as <Ctrl-A>.
		case lexer.TokenAngleLeft:
			p.Unscan()
			key, err := p.ParseSpecial()
			if err != nil {
				return nil, err
			}
			keyseq = append(keyseq, key)

		// Any other key than whitespace, end, or comment may be mapped for convenience.
		case lexer.TokenWhitespace, lexer.TokenEnd, lexer.TokenComment:
			p.Unscan()
			break Parse

		// Append to the key sequence list.
		default:
			seq := runeEventKeys(lit)
			keyseq = append(keyseq, seq...)
		}

		tok, lit = p.Scan()
	}

	if len(keyseq) == 0 {
		return nil, fmt.Errorf("Unexpected '%s', expected key sequence", lit)
	}

	return keyseq, nil
}

// Parse a special key name, such as <space> or <C-M-a>.
func (p *Parser) ParseSpecial() (key *tcell.EventKey, err error) {
	var mod tcell.ModMask

	// Scan the opening angle bracket
	tok, lit := p.Scan()
	if tok != lexer.TokenAngleLeft {
		return nil, fmt.Errorf("Unexpected %v, expected left angle bracket", lit)
	}

Scam:
	for {
		// Scan the next identifier, which must either be one of the modifiers
		// S, C, A, M, or an actual key name such as 'space' or 'f1'.
		tok, lit = p.Scan()
		if tok != lexer.TokenIdentifier {
			return nil, fmt.Errorf("Unexpected %v, expected identifier", lit)
		}

		// Turn key name into lowercase
		lit = strings.ToLower(lit)

		// Scan the next token, which may either be a sign, saying that
		// modifier keys are used, or a right angle bracket, which ends the
		// parsing with a key lookup. Any other key is an error.
		tok, _ := p.Scan()
		switch tok {
		case lexer.TokenAngleRight:
			break Scam
		case lexer.TokenPlus, lexer.TokenMinus:
			break
		default:
			return nil, fmt.Errorf("Unexpected '%s', expected >", lit)
		}

		// Apply the modifier key
		m, ok := modifiers[lit]
		if !ok {
			return nil, fmt.Errorf("Unexpected '%s', expected one of Shift, Ctrl, Alt, Meta", lit)
		}
		mod |= m
	}

	// Look up the parsed key name.
	ev, ok := keyNames[lit]
	if !ok {

		// If the key name is not found, it might be a letter.
		if len(lit) != 1 {
			return nil, fmt.Errorf("Unknown key name '%s'", lit)
		}

		// Return the rune and modifiers
		for _, r := range lit {
			return convertCtrlKey(tcell.NewEventKey(tcell.KeyRune, r, mod)), nil
		}
	}

	// Make a copy of the key, and apply any modifiers
	key = convertCtrlKey(tcell.NewEventKey(ev.Key(), ev.Rune(), mod))

	return key, nil
}

// runeEventKeys creates a slice of tcell.EventKey objects that corresponds to
// the literal runes in the string.
func runeEventKeys(s string) KeySequence {
	seq := make(KeySequence, 0)
	for _, r := range s {
		key := tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone)
		seq = append(seq, key)
	}
	return seq
}

// convertCtrlKey handles some special cases with Ctrl modifier keys, which are
// handled somewhat differently by tcell for a few specific cases.
func convertCtrlKey(ev *tcell.EventKey) *tcell.EventKey {
	modifiers := ev.Modifiers()
	hasCtrl := modifiers&tcell.ModCtrl == tcell.ModCtrl

	// If this is not a Ctrl+Rune event, return the original event.
	if !hasCtrl || ev.Key() != tcell.KeyRune {
		return ev
	}

	// Catch Ctrl+A through Ctrl+Z
	if ev.Rune() >= 'a' || ev.Rune() <= 'z' {
		ctrl := rune(tcell.KeyCtrlA) - 'a' + ev.Rune()
		return tcell.NewEventKey(tcell.Key(ctrl), rune(ctrl), modifiers)
	}

	// No more special rules
	return ev
}
