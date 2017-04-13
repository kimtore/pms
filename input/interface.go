package input

import (
	"fmt"

	"github.com/ambientsound/pms/input/commands"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/options"
)

type commandMap map[string]commands.Command

// Interface reads user input, tokenizes it, and dispatches the tokens to their respective commands.
type Interface struct {
	opts     *options.Options
	handlers commandMap
}

func NewInterface(opts *options.Options) *Interface {
	i := &Interface{opts: opts}
	i.handlers = make(commandMap, 0)
	return i
}

func (i *Interface) Execute(line string) error {
	var pos, nextPos int
	var token lexer.Token
	var cmd commands.Command
	var err error
	var ok bool

	for {
		token, nextPos = lexer.NextToken(line[pos:])
		pos += nextPos

		// First identifier; try to find a command handler
		if cmd == nil && token.Class == lexer.TokenIdentifier {
			key := token.String()
			if cmd, ok = i.handlers[key]; ok {
				cmd.Reset()
				continue
			}
			return fmt.Errorf("Not a command: %s", key)
		}

		err = cmd.Execute(token)

		if err != nil {
			return fmt.Errorf("Error while parsing input %s: %s", string(token.Runes), err)
		}
		if token.Class == lexer.TokenEnd {
			break
		}
	}

	return nil
}

func (i *Interface) Register(verb string, cmd commands.Command) error {
	if _, ok := i.handlers[verb]; ok {
		return fmt.Errorf("Handler with verb '%s' already exists", verb)
	}
	i.handlers[verb] = cmd
	return nil
}
