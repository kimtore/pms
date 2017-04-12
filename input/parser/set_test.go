package parser_test

import (
	"fmt"
	"testing"

	"github.com/ambientsound/pms/input/parser"
	"github.com/stretchr/testify/assert"
)

type parserTable struct {
	Name  string
	Error bool
	Input string
	Token parser.OptionToken
}

// TestOptionParser tests the variable tokenizer against a table of well-known inputs and outputs.
func TestOptionParser(t *testing.T) {
	table := []parserTable{
		{
			Name:  "string variable assignment",
			Input: "string=foo",
			Token: parser.OptionToken{
				Key:    "string",
				Value:  "foo",
				Bool:   false,
				Invert: false,
				Negate: false,
				Query:  false,
			},
		},
		{
			Name:  "variable query",
			Input: "string?",
			Token: parser.OptionToken{
				Key:    "string",
				Value:  "",
				Bool:   false,
				Invert: false,
				Negate: false,
				Query:  true,
			},
		},
		{
			Name:  "setting boolean option to true",
			Input: "bool",
			Token: parser.OptionToken{
				Key:    "bool",
				Value:  "",
				Bool:   true,
				Invert: false,
				Negate: false,
				Query:  false,
			},
		},
		{
			Name:  "setting boolean option to false",
			Input: "nobool",
			Token: parser.OptionToken{
				Key:    "bool",
				Value:  "",
				Bool:   true,
				Invert: false,
				Negate: true,
				Query:  false,
			},
		},
		{
			Name:  "inverting boolean option by 'inv' keyword",
			Input: "invbool",
			Token: parser.OptionToken{
				Key:    "bool",
				Value:  "",
				Bool:   true,
				Invert: true,
				Negate: false,
				Query:  false,
			},
		},
		{
			Name:  "inverting boolean option by exclamation mark",
			Input: "bool!",
			Token: parser.OptionToken{
				Key:    "bool",
				Value:  "",
				Bool:   true,
				Invert: true,
				Negate: false,
				Query:  false,
			},
		},
		{
			Name:  "negating boolean option starting with 'no'",
			Input: "nononsense",
			Token: parser.OptionToken{
				Key:    "nonsense",
				Value:  "",
				Bool:   true,
				Invert: false,
				Negate: true,
				Query:  false,
			},
		},
		{
			Name:  "querying boolean options while negating",
			Input: "noproblem?",
			Token: parser.OptionToken{
				Key:    "problem",
				Value:  "",
				Bool:   true,
				Invert: false,
				Negate: true,
				Query:  true,
			},
		},

		// Invalid queries
		{Input: "var!!", Error: true},
		{Input: "var!?", Error: true},
		{Input: "var?!", Error: true},
		{Input: "novar!?", Error: true},
		{Input: "novar?!", Error: true},
		{Input: "invvar!?", Error: true},
		{Input: "invvar?!", Error: true},
		{Input: "noinvvar", Error: true},
	}

	for i := range table {
		name := table[i].Name
		if len(name) == 0 {
			name = table[i].Input
		}
		input := table[i].Input
		check := table[i].Token
		token := parser.OptionToken{}
		err := token.Parse([]rune(input))
		if table[i].Error {
			assert.NotNil(t, err, fmt.Sprintf("Expected errors when parsing: %s", name))
		} else {
			assert.Nil(t, err, fmt.Sprintf("Expected no errors when parsing: %s", name))
			assert.Equal(t, check, token, fmt.Sprintf("Expected result when parsing: %s", name))
		}
	}
}
