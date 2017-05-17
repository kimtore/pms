package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/options"
)

// Set manipulates a Options table by parsing input tokens from the "set" command.
type Set struct {
	newcommand
	api    api.API
	tokens []parser.OptionToken
}

// NewSet returns Set.
func NewSet(api api.API) Command {
	return &Set{
		api:    api,
		tokens: make([]parser.OptionToken, 0),
	}
}

// Parse implements Command.
func (cmd *Set) Parse() error {

	for {
		// Scan the next token, which must be an identifier.
		tok, lit := cmd.ScanIgnoreWhitespace()
		switch tok {
		case lexer.TokenIdentifier:
			break
		case lexer.TokenEnd, lexer.TokenComment:
			return nil
		default:
			return fmt.Errorf("Unexpected '%s', expected whitespace or END", lit)
		}

		// Parse the option statement.
		cmd.Unscan()
		err := cmd.ParseSet()
		if err != nil {
			return err
		}
	}
}

// ParseSet parses a single "key=val" statement.
func (cmd *Set) ParseSet() error {
	tokens := make([]string, 0)
	for {
		tok, lit := cmd.Scan()
		if tok == lexer.TokenWhitespace || tok == lexer.TokenEnd || tok == lexer.TokenComment {
			break
		}
		tokens = append(tokens, lit)
	}

	s := strings.Join(tokens, "")
	optionToken := parser.OptionToken{}
	err := optionToken.Parse([]rune(s))
	if err != nil {
		return err
	}

	cmd.tokens = append(cmd.tokens, optionToken)

	return nil
}

// Exec implements Command.
func (cmd *Set) Exec() error {
	for _, tok := range cmd.tokens {
		opt := cmd.api.Options().Get(tok.Key)

		if opt == nil {
			return fmt.Errorf("No such option: %s", tok.Key)
		}

		// Queries print options to the statusbar.
		if tok.Query {
			cmd.api.Message(opt.String())
			continue
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

			// Not a boolean option, and no value. Print the value.
			cmd.api.Message(opt.String())
			continue
		}

		cmd.api.OptionChanged(opt.Key())
		cmd.api.Message(opt.String())
	}

	return nil
}
