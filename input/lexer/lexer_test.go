package lexer_test

import (
	"testing"

	"github.com/ambientsound/pms/input/lexer"
)

// TestLexer tests the lexer.NextToken() function, checking that it correctly
// splits up input lines into Token structs.
func TestLexer(t *testing.T) {
	i := 0
	pos := 0
	str := "the  quick <brown>\t \tfoxee # adds a comment"
	token := lexer.Token{}

	checks := []lexer.Token{
		lexer.Token{Class: lexer.TokenIdentifier, Runes: []rune("the")},
		lexer.Token{Class: lexer.TokenIdentifier, Runes: []rune("quick")},
		lexer.Token{Class: lexer.TokenIdentifier, Runes: []rune("<brown>")},
		lexer.Token{Class: lexer.TokenIdentifier, Runes: []rune("foxee")},
		lexer.Token{Class: lexer.TokenComment, Runes: []rune("# adds a comment")},
		lexer.Token{Class: lexer.TokenEnd, Runes: nil},
	}

	for {

		if i == len(checks) {
			if token.Class == lexer.TokenEnd {
				break
			}
			t.Fatalf("Tokenizer generated too many tokens!")
		}

		check := checks[i]
		token, npos := lexer.NextToken(str[pos:])
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
