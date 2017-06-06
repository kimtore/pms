package topbar

import (
	"fmt"
	"io"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/parser"
)

// Parser represents a parser.
type Parser struct {
	parser.Parser
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{parser.Parser{S: lexer.NewScanner(r)}}
}

// FragmentStatement holds information about a fragment, e.g.:
//
// ${variable|param}
// or:
// frag2
type FragmentStatement struct {
	Literal  string
	Variable string
	Param    string
	Instance Fragment
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

// ParseFragment parses a fragment statement.
func (p *Parser) ParseFragment() (*FragmentStatement, error) {
	stmt := &FragmentStatement{}

	tok, lit := p.Scan()

	switch tok {
	// If the token is a dollar sign, it refers to a variable, such as a tag in
	// the current song.
	case lexer.TokenVariable:
		break

	// If this is not a variable, use the text literally, and end parsing.
	default:
		stmt.Literal = lit
		return stmt, nil
	}

	// Next comes an identifier, or in case of parameterized variables, the
	// curly bracket opener.
	tok, lit = p.Scan()

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
	tok, lit = p.ScanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		return nil, fmt.Errorf("Unexpected %v, expected identifier", lit)
	}
	stmt.Variable = lit

	// Next, we can either have a separator in order to pass parameters, or
	// close with a curly bracket.
	tok, lit = p.ScanIgnoreWhitespace()

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
	tok, lit = p.ScanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		return nil, fmt.Errorf("Unexpected %v, expected parameter to $%s", lit, stmt.Variable)
	}
	stmt.Param = lit

	// The only valid token is the closing curly bracket.
	tok, lit = p.ScanIgnoreWhitespace()
	if tok != lexer.TokenClose {
		return nil, fmt.Errorf("Unexpected %v, expected '}'", lit)
	}

	return stmt, nil
}

// ParsePiece parses a piece statement.
func (p *Parser) ParsePiece() (*PieceStatement, error) {
	stmt := &PieceStatement{}

	for {
		tok, _ := p.Scan()

		// A piece is succeeded only by a new piece or new row.
		switch tok {
		case lexer.TokenStop:
			p.Unscan()
			fallthrough
		case lexer.TokenSeparator, lexer.TokenEnd:
			return stmt, nil
		}

		p.Unscan()
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
		tok, _ := p.Scan()

		// A row is succeeded only by a new row.
		switch tok {
		case lexer.TokenStop, lexer.TokenEnd:
			return stmt, nil
		}

		p.Unscan()
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
		tok, _ := p.Scan()

		// A matrix is never succeeded.
		switch tok {
		case lexer.TokenEnd:
			return stmt, nil
		}

		p.Unscan()
		row, err := p.ParseRow()
		if err != nil {
			return nil, err
		}

		stmt.Rows = append(stmt.Rows, row)
	}
}
