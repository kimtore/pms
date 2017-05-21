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

func (p *Parser) ParseKeySequence() (tok int, lit string, seq KeySequence, err error) {
	for {
		tok, lit = p.Scan()
	}
}

// ParseKey parses the next literal key.
func (p *Parser) ParseKey() (tok int, lit string, key tcell.EventKey, err error) {
	tok, lit = p.Scan()
	if tok == lexer.TokenWhitespace {
		return
	}
	if tok == lexer.TokenAngleLeft {
		p.Unscan()
		return p.ParseSpecial()
	}
	return
}

// Parse a special key name, such as <space> or <C-M-a>.
func (p *Parser) ParseSpecial() (tok int, lit string, key tcell.EventKey, err error) {
	var mod tcell.ModMask

	// Scan the opening angle bracket
	tok, lit = p.Scan()
	if tok != lexer.TokenAngleLeft {
		err = fmt.Errorf("Unexpected %v, expected left angle bracket", lit)
		return
	}

Scam:
	for {
		// Scan the next identifier, which must either be one of the modifiers
		// S, C, A, M, or an actual key name such as 'space' or 'f1'.
		tok, lit = p.Scan()
		if tok != lexer.TokenIdentifier {
			err = fmt.Errorf("Unexpected %v, expected identifier", lit)
			return
		}

		// Scan the next token, which may either be a sign, saying that
		// modifier keys are used, or a right angle bracket, which ends the
		// parsing with a key lookup. Any other key is an error.
		k := lit
		tok, lit = p.Scan()
		switch tok {
		case lexer.TokenAngleRight:
			break Scam
		case lexer.TokenPlus, lexer.TokenMinus:
			lit = k
			break
		default:
			err = fmt.Errorf("Unexpected %v, expected >", lit)
			return
		}

		// Apply the modifier key
		m, ok := modifiers[strings.ToLower(lit)]
		if !ok {
			err = fmt.Errorf("Unexpected %v, expected one of Shift, Ctrl, Alt, Meta", lit)
			return
		}
		mod |= m
	}

	// Look up the parsed key name.
	ev, ok := keyNames[lit]
	if !ok {
		err = fmt.Errorf("Unknown key name %v", lit)
		return
	}

	tok = lexer.TokenIdentifier
	key = *ev

	return
}
