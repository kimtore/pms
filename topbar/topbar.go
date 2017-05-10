package topbar

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
)

// Fragment is the smallest possible unit in a topbar.
type Fragment interface {

	// Text returns a string that should be drawn,
	// along with its stylesheet identifier.
	Text() (string, string)
}

// fragments is a map of fragments that can be drawn in the topbar, along with
// their textual representation. When implementing a new topbar fragment, place
// its constructor in this map.
var fragments = map[string]func(api.API, string) Fragment{
	"mode":      NewMode,
	"shortname": NewShortname,
	"tag":       NewTag,
	"version":   NewVersion,
	"volume":    NewVolume,
}

// NewFragment constructs a new Fragment based on a parsed topbar fragment statement.
func NewFragment(a api.API, stmt *FragmentStatement) (Fragment, error) {
	if len(stmt.Variable) == 0 {
		return NewText(stmt.Literal), nil
	}
	if ctor, ok := fragments[stmt.Variable]; ok {
		return ctor(a, stmt.Param), nil
	}
	return nil, fmt.Errorf("Unrecognized variable '${%s}'", stmt.Variable)
}

// Parse sets up a lexer and parser for a topbar matrix statement, instantiates
// fragments, and returns the parse tree.
func Parse(a api.API, input string) (*MatrixStatement, error) {
	reader := strings.NewReader(input)
	parser := NewParser(reader)

	matrixStmt, err := parser.ParseMatrix()
	if err != nil {
		return nil, err
	}

	// Instantiate fragments
	for _, rowStmt := range matrixStmt.Rows {
		for _, pieceStmt := range rowStmt.Pieces {
			for _, fragmentStmt := range pieceStmt.Fragments {
				frag, err := NewFragment(a, fragmentStmt)
				if err != nil {
					return nil, err
				}
				fragmentStmt.Instance = frag
			}
		}
	}

	return matrixStmt, nil
}
