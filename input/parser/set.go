package parser

import (
	"fmt"
)

// OptionToken represents a key=value string, which can have setting, inversion, negation, and queries.
type OptionToken struct {
	Key    string
	Value  string
	Bool   bool
	Negate bool
	Invert bool
	Query  bool
}

// Parse parses a option=value string.
func (t *OptionToken) Parse(runes []rune) error {
	// Parsing the value is done verbatim, whereas the key has
	// modifiers such as !, ?, inv*, no*.
	parsingKey := true

	for _, r := range runes {
		if !parsingKey {
			t.Value += string(r)
			continue
		}

		if t.Query {
			return fmt.Errorf("Trailing characters after '?'")
		} else if r == '=' {
			parsingKey = false
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
	if parsingKey && !t.Query {
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
