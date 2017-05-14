package parser

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
)

// Parser represents a parser.
type Parser struct {
	S   *lexer.Scanner
	buf struct {
		tok int    // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// Scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) Scan() (tok int, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.S.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// ScanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) ScanIgnoreWhitespace() (tok int, lit string) {
	tok, lit = p.Scan()
	if tok == lexer.TokenWhitespace {
		tok, lit = p.Scan()
	}
	return
}

// Unscan pushes the previously read token back onto the buffer.
func (p *Parser) Unscan() { p.buf.n = 1 }

// ParseEnd parses to the end, and returns an error if the end hasn't been reached.
func (p *Parser) ParseEnd() error {
	tok, lit := p.ScanIgnoreWhitespace()
	if tok != lexer.TokenEnd {
		return fmt.Errorf("Unexpected %v, expected END", lit)
	}
	return nil
}
