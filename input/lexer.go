package input

import "unicode"

type Token struct {
	Class string
	Runes []rune
}

func NewToken() Token {
	t := Token{}
	t.Runes = make([]rune, 0)
	return t
}

// Tokenize is a lexer for the PMS input language. It splits up an input sentence into tokens, separating
func Tokenize(s string) []Token {
	tokens := make([]Token, 0)
	t := NewToken()
	for _, r := range s {
		if t.Class != "comment" && unicode.IsSpace(r) {
			if len(t.Runes) > 0 {
				tokens = append(tokens, t)
				t = NewToken()
			}
			continue
		} else if len(t.Runes) == 0 {
			switch r {
			case '#':
				// Comments terminate the line
				t.Class = "comment"
			default:
				t.Class = "identifier"
			}
		}
		t.Runes = append(t.Runes, r)
	}
	if len(t.Runes) > 0 {
		tokens = append(tokens, t)
	}
	return tokens
}
