package lexer_test

import (
	"testing"

	"github.com/ambientsound/pms/input/lexer"
)

type result struct {
	class int
	str   string
}

var lexerTests = []struct {
	input    string
	expected []result
}{
	{
		"the  quick <brown>\t \tfoxee # adds a comment",
		[]result{
			{class: lexer.TokenIdentifier, str: "the"},
			{class: lexer.TokenIdentifier, str: "quick"},
			{class: lexer.TokenIdentifier, str: "<brown>"},
			{class: lexer.TokenIdentifier, str: "foxee"},
			{class: lexer.TokenComment, str: "# adds a comment"},
			{class: lexer.TokenEnd, str: ""},
		},
	},
	{
		"$variables are {nice }}, ar{}en't $they?",
		[]result{
			{class: lexer.TokenVariable, str: "$"},
			{class: lexer.TokenIdentifier, str: "variables"},
			{class: lexer.TokenIdentifier, str: "are"},
			{class: lexer.TokenOpen, str: "{"},
			{class: lexer.TokenIdentifier, str: "nice"},
			{class: lexer.TokenClose, str: "}"},
			{class: lexer.TokenClose, str: "}"},
			{class: lexer.TokenIdentifier, str: ","},
			{class: lexer.TokenIdentifier, str: "ar"},
			{class: lexer.TokenOpen, str: "{"},
			{class: lexer.TokenClose, str: "}"},
			{class: lexer.TokenIdentifier, str: "en't"},
			{class: lexer.TokenVariable, str: "$"},
			{class: lexer.TokenIdentifier, str: "they?"},
			{class: lexer.TokenEnd, str: ""},
		},
	},
	{
		"$1$2",
		[]result{
			{class: lexer.TokenVariable, str: "$"},
			{class: lexer.TokenIdentifier, str: "1"},
			{class: lexer.TokenVariable, str: "$"},
			{class: lexer.TokenIdentifier, str: "2"},
			{class: lexer.TokenEnd, str: ""},
		},
	},
}

// TestLexer tests the lexer.NextToken() function, checking that it correctly
// splits up input lines into Token structs.
func TestLexer(t *testing.T) {
	var token lexer.Token

	for _, test := range lexerTests {

		i := 0
		pos := 0

		for {

			if i == len(test.expected) {
				if token.Class == lexer.TokenEnd {
					break
				}
				t.Fatalf("Tokenizer generated too many tokens!")
			}

			check := test.expected[i]
			token, npos := lexer.NextToken(test.input[pos:])
			pos += npos
			str := token.String()

			t.Logf("Token %d: pos=%d, runes='%s', input='%s'", i, pos, str, test.input)

			if token.Class != check.class {
				t.Fatalf("Token class for token %d is wrong; expected %d but got %d", i+1, check.class, token.Class)
			}

			if check.str != str {
				t.Fatalf("String check against token %d failed; expected '%s' but got '%s'", i+1,
					check.str,
					str,
				)
			}

			i++
		}
	}
}
