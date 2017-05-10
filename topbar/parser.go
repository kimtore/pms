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

// PieceStatement holds information about a piece, e.g.:
//
// ${variable|param} frag2
type PieceStatement struct {
	Fragments []*FragmentStatement
}

// RowStatement holds information about a row, e.g.:
//
// ${variable|param} frag2|text2|text3
type RowStatement struct {
	Pieces []*PieceStatement
}

// MatrixStatement is an initialization of a complete topbar, e.g.:
//
// ${variable|param} frag2|text2|text3;row1.1|row1.2|row1.3
type MatrixStatement struct {
	Rows []*RowStatement
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
	stmt := &FragmentStatement{}

	tok, lit := p.scan()

	switch tok {
	// The first token should either be whitespace or an identifier, which
	// will compose the entirety of the fragment.
	case lexer.TokenWhitespace, lexer.TokenIdentifier:
		stmt.Literal = lit
		return stmt, nil

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
		stmt.Variable = lit
		return stmt, nil

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
	stmt.Variable = lit

	// Next, we can either have a separator in order to pass parameters, or
	// close with a curly bracket.
	tok, lit = p.scanIgnoreWhitespace()

	switch tok {
	// Parameterized variable, e.g. '${tag|artist}'.
	case lexer.TokenSeparator:
		break

	// Finished parsing the variable.
	case lexer.TokenClose:
		return stmt, nil

	// No other tokens are valid.
	default:
		return nil, fmt.Errorf("Unexpected %v, expected '|' or '}'", lit)
	}

	// The only valid token is the parameter.
	tok, lit = p.scanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		return nil, fmt.Errorf("Unexpected %v, expected parameter to $%s", lit, stmt.Variable)
	}
	stmt.Param = lit

	// The only valid token is the closing curly bracket.
	tok, lit = p.scanIgnoreWhitespace()
	if tok != lexer.TokenClose {
		return nil, fmt.Errorf("Unexpected %v, expected '}'", lit)
	}

	return stmt, nil
}

// ParsePiece parses a piece statement.
func (p *Parser) ParsePiece() (*PieceStatement, error) {
	stmt := &PieceStatement{}

	for {
		tok, _ := p.scan()

		// A piece is succeeded only by a new piece or new row.
		switch tok {
		case lexer.TokenStop:
			p.unscan()
			fallthrough
		case lexer.TokenSeparator, lexer.TokenEnd:
			return stmt, nil
		}

		p.unscan()
		frag, err := p.ParseFragment()
		if err != nil {
			return nil, err
		}

		stmt.Fragments = append(stmt.Fragments, frag)
	}
}

// ParsePiece parses a row statement.
func (p *Parser) ParseRow() (*RowStatement, error) {
	stmt := &RowStatement{}

	for {
		tok, _ := p.scan()

		// A row is succeeded only by a new row.
		switch tok {
		case lexer.TokenStop, lexer.TokenEnd:
			return stmt, nil
		}

		p.unscan()
		piece, err := p.ParsePiece()
		if err != nil {
			return nil, err
		}

		stmt.Pieces = append(stmt.Pieces, piece)
	}
}

// ParseMatrix parses a matrix statement.
func (p *Parser) ParseMatrix() (*MatrixStatement, error) {
	stmt := &MatrixStatement{}

	for {
		tok, _ := p.scan()

		// A matrix is never succeeded.
		switch tok {
		case lexer.TokenEnd:
			return stmt, nil
		}

		p.unscan()
		row, err := p.ParseRow()
		if err != nil {
			return nil, err
		}

		stmt.Rows = append(stmt.Rows, row)
	}
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
