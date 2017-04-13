package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/options"
)

// Set manipulates a Options table by parsing input tokens from the "set" command.
type Set struct {
	opts   *options.Options
	tokens []parser.OptionToken
}

func NewSet(opts *options.Options) *Set {
	p := &Set{}
	p.opts = opts
	p.tokens = make([]parser.OptionToken, 0)
	return p
}

func (p *Set) Parse(t input.Token) error {
	if t.Class == input.TokenEnd {
		return nil
	}
	if t.Class != input.TokenIdentifier {
		return fmt.Errorf("Unknown input '%s', expected identifier", string(t.Runes))
	}
	tok := parser.OptionToken{}
	err := tok.Parse(t.Runes)
	if err != nil {
		return err
	}

	opt := p.opts.Get(tok.Key)

	if tok.Query {
		return nil // FIXME: statusbar feedback
	}

	switch opt := opt.(type) {
	case *options.BoolOption:
		if !tok.Bool {
			return fmt.Errorf("Attempting to give parameters to a boolean option (try 'set (no|inv)?%s')", tok.Key)
		}
		if tok.Invert {
			opt.SetBool(!opt.BoolValue())
			return nil
		}
		if tok.Negate {
			opt.SetBool(false)
			return nil
		}
		opt.SetBool(true)
		return nil
	default:
		if !tok.Bool {
			return opt.Set(tok.Value)
		}
		return fmt.Errorf("Attempting to execute a boolean operation on a non-boolean option")
	}

	return nil
}
