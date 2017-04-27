package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/options"
)

// Set manipulates a Options table by parsing input tokens from the "set" command.
type Set struct {
	opts     *options.Options
	messages chan string
}

func NewSet(opts *options.Options, messages chan string) *Set {
	p := &Set{}
	p.opts = opts
	p.messages = messages
	p.Reset()
	return p
}

func (p *Set) Reset() {
}

func (p *Set) Execute(t lexer.Token) error {
	s := t.String()

	if t.Class == lexer.TokenEnd {
		return nil
	}
	if t.Class != lexer.TokenIdentifier {
		return fmt.Errorf("Unknown input '%s', expected identifier", s)
	}
	tok := parser.OptionToken{}
	err := tok.Parse(t.Runes)
	if err != nil {
		return err
	}

	opt := p.opts.Get(tok.Key)

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
			p.message(opt)
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

msg:
	p.message(opt)
	return nil
}

func (cmd *Set) message(opt options.Option) {
	cmd.messages <- opt.String()
}
