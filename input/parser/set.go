package parser

import (
	"fmt"

	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/options"
)

type OptionToken struct {
	Key    string
	Value  string
	Bool   bool
	Negate bool
	Invert bool
	Query  bool
}

// SetParser parses input text starting with "set"
type SetParser struct {
	opts   *options.Options
	tokens []OptionToken
}

func NewSetParser() *SetParser {
	p := &SetParser{}
	p.tokens = make([]OptionToken, 0)
	return p
}

// Parse parses a option=value string, accounting for inversion, negation, and queries.
func (t *OptionToken) Parse(runes []rune) error {
	// Parsing the value is done verbatim, whereas the key has
	// modifiers such as !, ?, inv*, no*.
	parsing_key := true

	for _, r := range runes {
		if !parsing_key {
			t.Value += string(r)
			continue
		}

		if t.Query {
			return fmt.Errorf("Trailing characters after '?'")
		} else if r == '=' {
			parsing_key = false
		} else if r == '?' {
			t.Query = true
		} else if r == '!' {
			if t.Invert {
				return fmt.Errorf("Double inversion not allowed")
			}
			t.Invert = true
		} else {
			t.Key += string(r)
			if t.Key == "no" && !t.Negate {
				t.Key = ""
				t.Negate = true
				t.Bool = true
			} else if t.Key == "inv" && !t.Invert {
				t.Key = ""
				t.Invert = true
				t.Bool = true
			}
		}
	}
	if parsing_key && !t.Query {
		t.Bool = true
	}
	if t.Query {
		if t.Invert {
			return fmt.Errorf("Query operation cannot be combined with inversion")
		}
	} else {
		if t.Negate && t.Invert {
			return fmt.Errorf("Negation and inversion cannot be combined")
		}
	}
	return nil
}

func (p *SetParser) Parse(t input.Token) error {
	if t.Class == input.TokenEnd {
		return nil
	}
	if t.Class != input.TokenIdentifier {
		return fmt.Errorf("Unknown input '%s', expected identifier", string(t.Runes))
	}
	tok := OptionToken{}
	err := tok.Parse(t.Runes)
	if err != nil {
		p.tokens = append(p.tokens, tok)
	}
	return err
}
