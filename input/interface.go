package input

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/input/lexer"
)

// CLI reads user input, tokenizes it, and dispatches the tokens to their respective commands.
type CLI struct {
	baseAPI api.API
}

func NewCLI(baseAPI api.API) *CLI {
	return &CLI{
		baseAPI: baseAPI,
	}
}

func (i *CLI) Execute(line string) error {
	var cmd commands.Command
	var err error

	reader := strings.NewReader(line)
	scanner := lexer.NewScanner(reader)

	for {
		class, token := scanner.Scan()

		// First identifier; try to find a command handler
		if cmd == nil {
			switch class {
			case lexer.TokenIdentifier:
				if ctor, ok := commands.Verbs[token]; ok {
					cmd = ctor(i.baseAPI)
					continue
				}
				return fmt.Errorf("Not a command: %s", token)
			case lexer.TokenComment:
				continue
			case lexer.TokenEnd:
				return nil
			case lexer.TokenStop:
				cmd = nil
				continue
			case lexer.TokenWhitespace:
				continue
			default:
				return fmt.Errorf("Unexpected '%s', expected identifier", token)
			}
		}

		if class == lexer.TokenWhitespace {
			continue
		}

		err = cmd.Execute(class, token)

		if err != nil {
			return err
		}

		if class == lexer.TokenEnd {
			break
		}
	}

	return nil
}
