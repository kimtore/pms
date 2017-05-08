package lexer

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
	TokenWhitespace
	TokenIdentifier
	TokenComment
	TokenVariable
	TokenOpen
	TokenClose
)

// runeClass returns the token class of an input character.
func runeClass(r rune) int {
	if unicode.IsSpace(r) {
		return TokenWhitespace
	}
	switch r {
	case '$':
		return TokenVariable
	case '{':
		return TokenOpen
	case '}':
		return TokenClose
	case '#':
		return TokenComment
	default:
		return TokenIdentifier
	}
}

// Tokenize is a lexer for the input language. It extracts the next token out
// of a sentence, and returns the token itself, in addition to the next position
// in the string that should be parsed.
func NextToken(input string) (t Token, pos int) {
	t = NewToken()

	for _, r := range input {
		firstChar := len(t.Runes) == 0
		runeLen := utf8.RuneLen(r)
		class := runeClass(r)

		if t.Class == TokenComment {
			t.Runes = append(t.Runes, r)
			pos += runeLen
			continue
		}

		if class == TokenWhitespace {
			pos += runeLen
			if len(t.Runes) > 0 {
				return
			}
			continue
		}

		if firstChar {
			t.Class = class
			t.Runes = append(t.Runes, r)
			pos += runeLen
			continue
		}

		if class == TokenIdentifier && t.Class == class {
			t.Runes = append(t.Runes, r)
			pos += runeLen
			continue
		}

		break
	}

	return
}

func (t *Token) String() string {
	return string(t.Runes)
}
