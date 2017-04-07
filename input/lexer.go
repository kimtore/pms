package input

import (
	"unicode"
	"unicode/utf8"
)

type Token struct {
	Class int
	Runes []rune
}

func NewToken() Token {
	t := Token{}
	t.Runes = make([]rune, 0)
	return t
}

const (
	TokenEnd = iota
	TokenIdentifier
	TokenComment
)

// Tokenize is a lexer for the input language. It extracts the next token out of a sentence, and returns
func NextToken(input string) (t Token, pos int) {
	t = NewToken()
	for _, r := range input {
		pos += utf8.RuneLen(r)
		if t.Class != TokenComment && unicode.IsSpace(r) {
			if len(t.Runes) > 0 {
				return
			}
			continue
		} else if len(t.Runes) == 0 {
			switch r {
			case '#':
				// Comments terminate the line
				t.Class = TokenComment
			default:
				t.Class = TokenIdentifier
			}
		}
		t.Runes = append(t.Runes, r)
	}
	return
}
