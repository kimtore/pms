package actions

import (
	"fmt"

	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/options"
)

// SetParser parses input text starting with "set"
type SetParser struct {
	opts   *options.Options
	tokens []parser.OptionToken
}

func NewSetParser() *SetParser {
	p := &SetParser{}
	p.tokens = make([]parser.OptionToken, 0)
	return p
}

func (p *SetParser) Parse(t input.Token) error {
	if t.Class == input.TokenEnd {
		return nil
	}
	if t.Class != input.TokenIdentifier {
		return fmt.Errorf("Unknown input '%s', expected identifier", string(t.Runes))
	}
	tok := parser.OptionToken{}
	err := tok.Parse(t.Runes)
	if err != nil {
		p.tokens = append(p.tokens, tok)
	}
	return err
}
