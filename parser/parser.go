package parser

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
)

// Token represent a token: classification and literal text
type Token struct {
	Tok int
	Lit string
}

// buf represents the last scanned token.
type buf struct {
	Token
	n int // buffer size (max=1)
}

// Parser represents a parser.
type Parser struct {
	S       *lexer.Scanner
	buf     buf
	scanned []Token
}

// New returns Parser.
func New(r *lexer.Scanner) *Parser {
	return &Parser{S: r}
}

// SetScanner assigns a scanner object to the parser.
func (p *Parser) SetScanner(s *lexer.Scanner) {
	p.S = s
	p.buf.n = 0
}

// Scanned returns all scanned tokens.
func (p *Parser) Scanned() []Token {
	return p.scanned
}

// Scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) Scan() (tok int, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.Tok, p.buf.Lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.S.Scan()

	// Create the scanned buffer.
	if p.scanned == nil {
		p.scanned = make([]Token, 0)
	}

	// Push the data to the scanned buffer.
	p.scanned = append(p.scanned, Token{tok, lit})

	// Save it to the buffer in case we unscan later.
	p.buf.Tok, p.buf.Lit = tok, lit

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
