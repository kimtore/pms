package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/options"
)

// Set manipulates a Options table by parsing input tokens from the "set" command.
type Set struct {
	api API
}

func NewSet(api API) Command {
	return &Set{
		api: api,
	}
}

func (p *Set) Execute(t lexer.Token) error {
	if t.Class == lexer.TokenEnd || t.Class == lexer.TokenComment {
		return nil
	}

	tok := parser.OptionToken{}
	err := tok.Parse(t.Runes)
	if err != nil {
		return err
	}

	opt := p.api.Options().Get(tok.Key)

	if opt == nil {
		return fmt.Errorf("No such option: %s", tok.Key)
	}

	if tok.Query {
		goto msg
	}

	switch opt := opt.(type) {

	case *options.BoolOption:
		switch {
		case !tok.Bool:
			return fmt.Errorf("Attempting to give parameters to a boolean option (try 'set no%s' or 'set inv%s')", tok.Key, tok.Key)
		case tok.Invert:
			opt.SetBool(!opt.BoolValue())
			p.api.Message(opt.String())
		case tok.Negate:
			opt.SetBool(false)
		default:
			opt.SetBool(true)
		}

	default:
		if !tok.Bool {
			if err := opt.Set(tok.Value); err != nil {
				return err
			}
			break
		}
		return fmt.Errorf("Attempting to execute a boolean operation on a non-boolean option")
	}

	p.api.OptionChanged(opt.Key())

msg:
	p.api.Message(opt.String())
	return nil
}
