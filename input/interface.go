package input

import (
	"fmt"

	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/input/lexer"
)

type commandCtor func(commands.API) commands.Command

type commandMap map[string]commandCtor

// CLI reads user input, tokenizes it, and dispatches the tokens to their respective commands.
type CLI struct {
	handlers commandMap
	baseAPI  commands.API
}

func NewCLI(baseAPI commands.API) *CLI {
	return &CLI{
		baseAPI:  baseAPI,
		handlers: make(commandMap, 0),
	}
}

func (i *CLI) Execute(line string) error {
	var pos, nextPos int
	var token lexer.Token
	var cmd commands.Command
	var err error

	for {
		token, nextPos = lexer.NextToken(line[pos:])
		pos += nextPos

		// First identifier; try to find a command handler
		if cmd == nil {
			key := token.String()
			switch token.Class {
			case lexer.TokenIdentifier:
				if ctor, ok := i.handlers[key]; ok {
					cmd = ctor(i.baseAPI)
					continue
				}
				return fmt.Errorf("Not a command: %s", key)
			case lexer.TokenComment:
				continue
			case lexer.TokenEnd:
				return nil
			default:
				return fmt.Errorf("Unexpected '%s', expected identifier", key)
			}
		}

		err = cmd.Execute(token)

		if err != nil {
			return err
		}

		if token.Class == lexer.TokenEnd {
			break
		}
	}

	return nil
}

func (i *CLI) Register(verb string, ctor commandCtor) error {
	if _, ok := i.handlers[verb]; ok {
		return fmt.Errorf("Handler with verb '%s' already exists", verb)
	}
	i.handlers[verb] = ctor
	return nil
}
