package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/options"
)

// Set manipulates a Options table by parsing input tokens from the "set" command.
type Set struct {
	api   api.API
	key   string
	val   string
	equal bool
}

func NewSet(api api.API) Command {
	return &Set{
		api: api,
	}
}

func (cmd *Set) Execute(class int, s string) error {

	switch class {

	case lexer.TokenIdentifier:
		if len(cmd.key) > 0 {
			if len(cmd.val) > 0 {
				return fmt.Errorf("Unexpected '%s', expected whitespace or END", s)
			}
			cmd.val = s
		} else {
			cmd.key = s
		}

	case lexer.TokenEqual:
		if len(cmd.key) == 0 {
			return fmt.Errorf("Unexpected '%s', expected option", s)
		}
		if cmd.equal {
			return fmt.Errorf("Unexpected '%s', expected option value", s)
		}
		cmd.equal = true

	case lexer.TokenEnd, lexer.TokenComment, lexer.TokenWhitespace:
		if len(cmd.key) == 0 {
			cmd.reset()
			return nil
		}
		token := cmd.key
		if len(cmd.val) > 0 {
			token = fmt.Sprintf("%s=%s", cmd.key, cmd.val)
		}
		cmd.reset()
		return cmd.run(token)

	default:
		return fmt.Errorf("Unexpected '%s', expected identifier", s)
	}

	return nil
}

// reset resets internal state between options
func (cmd *Set) reset() {
	cmd.key = ""
	cmd.val = ""
	cmd.equal = false
}

func (cmd *Set) run(token string) error {

	tok := parser.OptionToken{}
	err := tok.Parse([]rune(token))
	if err != nil {
		return err
	}

	opt := cmd.api.Options().Get(tok.Key)

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
			cmd.api.Message(opt.String())
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
		return fmt.Errorf("Attempting to set '%s', but '%s' is not a boolean option", token, tok.Key)
	}

	cmd.api.OptionChanged(opt.Key())

msg:
	cmd.api.Message(opt.String())
	return nil
}
