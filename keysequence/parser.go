package keysequence

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/parser"
	"github.com/ambientsound/pms/term"
	termbox "github.com/nsf/termbox-go"
)

// modifiers maps modifier names to their integer representation.
var modifiers = map[string]term.Modifier{
	"s":     term.ModShift,
	"c":     term.ModCtrl,
	"a":     term.ModAlt,
	"m":     term.ModMeta,
	"shift": term.ModShift,
	"ctrl":  term.ModCtrl,
	"alt":   term.ModAlt,
	"meta":  term.ModMeta,
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
func (p *Parser) ParseSpecial() (key term.KeyPress, err error) {
	var mod term.Modifier

	// Scan the opening angle bracket
	tok, lit := p.Scan()
	if tok != lexer.TokenAngleLeft {
		return key, fmt.Errorf("Unexpected %v, expected left angle bracket", lit)
	}

Scam:
	for {
		// Scan the next identifier, which must either be one of the modifiers
		// S, C, A, M, or an actual key name such as 'space' or 'f1'.
		tok, lit = p.Scan()
		if tok != lexer.TokenIdentifier {
			return key, fmt.Errorf("Unexpected %v, expected identifier", lit)
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
			return key, fmt.Errorf("Unexpected '%s', expected >", lit)
		}

		// Apply the modifier key
		m, ok := modifiers[lit]
		if !ok {
			return key, fmt.Errorf("Unexpected '%s', expected one of Shift, Ctrl, Alt, Meta", lit)
		}
		mod |= m
	}

	// Look up the last part of the parsed key name.
	ev, err := term.Key(lit)
	if err == nil {
		key = term.KeyPress{ev, 0, mod}
		if ev == termbox.KeySpace {
			key.Ch = ' '
		}
		return key, nil
	}

	// If the key name is not found, it is either incorrect, or a letter.
	if len(lit) != 1 {
		return key, fmt.Errorf("Unknown key name '%s'", lit)
	}

	// Return the rune and modifiers
	for _, r := range lit {
		key := termbox.Key(0)
		if mod == term.ModCtrl {
			key = termbox.Key(r - 96)
		}
		return term.KeyPress{key, r, mod}, nil
	}

	return key, nil
}

// runeEventKeys creates a slice ofterm.KeyPress objects that corresponds to
// the literal runes in the string.
func runeEventKeys(s string) KeySequence {
	seq := make(KeySequence, 0)
	for _, r := range s {
		key := term.KeyPress{0, r, 0}
		seq = append(seq, key)
	}
	return seq
}
