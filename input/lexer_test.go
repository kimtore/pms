package input

import (
	"testing"

	"github.com/ambientsound/pms/input"
)

func TestLexer(t *testing.T) {
	str := "the  quick <brown>\t \tfoxee # adds a comment"
	check := []input.Token{
		input.Token{Runes: []rune("the"), Class: "identifier"},
		input.Token{Runes: []rune("quick"), Class: "identifier"},
		input.Token{Runes: []rune("<brown>"), Class: "identifier"},
		input.Token{Runes: []rune("foxee"), Class: "identifier"},
		input.Token{Runes: []rune("# adds a comment"), Class: "comment"},
	}
	tokens := input.Tokenize(str)
	if len(check) != len(tokens) {
		t.FailNow()
	}
	for i := range tokens {
		if tokens[i].Class != check[i].Class {
			t.FailNow()
		}
		if string(check[i].Runes) != string(tokens[i].Runes) {
			t.FailNow()
		}
	}
}
