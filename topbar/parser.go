package topbar

import (
	"fmt"
	"io"

	"github.com/ambientsound/pms/input/lexer"
)

// FragmentStatement holds information about a fragment, e.g.:
//
// ${variable|param}
// or:
// frag2
type FragmentStatement struct {
	Literal  string
	Variable string
	Param    string
}

// Parser represents a parser.
type Parser struct {
	s   *lexer.Scanner
	buf struct {
		tok int    // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: lexer.NewScanner(r)}
}

// ParseFragment parses a fragment statement.
func (p *Parser) ParseFragment() (*FragmentStatement, error) {
	frag := &FragmentStatement{}

	tok, lit := p.scan()

	switch tok {
	// The first token should either be whitespace or an identifier, which
	// will compose the entirety of the fragment.
	case lexer.TokenWhitespace, lexer.TokenIdentifier:
		frag.Literal = lit
		return frag, nil

	// Otherwise, look for a dollar sign, which refers to a dynamic topbar
	// fragment, such as a tag in the current song.
	case lexer.TokenVariable:
		break

	// No other tokens are valid.
	default:
		return nil, fmt.Errorf("Unexpected %v, expected '$'", lit)
	}

	// Next comes an identifier, or in case of parameterized variables, the
	// curly bracket opener.
	tok, lit = p.scan()

	switch tok {
	// Simple variable, e.g. '$artist'. Return immediately.
	case lexer.TokenIdentifier:
		frag.Variable = lit
		return frag, nil

	// Curly bracket opener, e.g. '${tag}'.
	case lexer.TokenOpen:
		break

	// No other tokens are valid.
	default:
		return nil, fmt.Errorf("Unexpected %v, expected '{' or identifier", lit)
	}

	// Next, the only valid token is the variable name.
	tok, lit = p.scanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		return nil, fmt.Errorf("Unexpected %v, expected identifier", lit)
	}
	frag.Variable = lit

	// Next, we can either have a separator in order to pass parameters, or
	// close with a curly bracket.
	tok, lit = p.scanIgnoreWhitespace()

	switch tok {
	// Parameterized variable, e.g. '${tag|artist}'.
	case lexer.TokenSeparator:
		break

	// Finished parsing the variable.
	case lexer.TokenClose:
		return frag, nil

	// No other tokens are valid.
	default:
		return nil, fmt.Errorf("Unexpected %v, expected '|' or '}'", lit)
	}

	// The only valid token is the parameter.
	tok, lit = p.scanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		return nil, fmt.Errorf("Unexpected %v, expected parameter to $%s", lit, frag.Variable)
	}
	frag.Param = lit

	// The only valid token is the closing curly bracket.
	tok, lit = p.scanIgnoreWhitespace()
	if tok != lexer.TokenClose {
		return nil, fmt.Errorf("Unexpected %v, expected '}'", lit)
	}

	return frag, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok int, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok int, lit string) {
	tok, lit = p.scan()
	if tok == lexer.TokenWhitespace {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
