package input_test

import (
	"testing"

	"github.com/ambientsound/pms/input"
)

// TestLexer tests the input.NextToken() function, checking that it correctly
// splits up input lines into Token structs.
func TestLexer(t *testing.T) {
	i := 0
	pos := 0
	str := "the  quick <brown>\t \tfoxee # adds a comment"
	token := input.Token{}

	checks := []input.Token{
		input.Token{Class: input.TokenIdentifier, Runes: []rune("the")},
		input.Token{Class: input.TokenIdentifier, Runes: []rune("quick")},
		input.Token{Class: input.TokenIdentifier, Runes: []rune("<brown>")},
		input.Token{Class: input.TokenIdentifier, Runes: []rune("foxee")},
		input.Token{Class: input.TokenComment, Runes: []rune("# adds a comment")},
		input.Token{Class: input.TokenEnd, Runes: nil},
	}

	for {

		if i == len(checks) {
			if token.Class == input.TokenEnd {
				break
			}
			t.Fatalf("Tokenizer generated too many tokens!")
		}

		check := checks[i]
		token, npos := input.NextToken(str[pos:])
		pos += npos

		t.Logf("Token %d: pos=%d, runes='%s', input='%s'", i, pos, string(token.Runes), str)

		if token.Class != check.Class {
			t.Fatalf("Token class for token %d is wrong; expected %d but got %d", i+1, check.Class, token.Class)
		}

		if string(check.Runes) != string(token.Runes) {
			t.Fatalf("String check against token %d failed; expected '%s' but got '%s'", i+1,
				string(check.Runes),
				string(token.Runes),
			)
		}

		i++
	}
}
