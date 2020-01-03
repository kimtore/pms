package input

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/input/lexer"
)

// Interpreter reads user input, tokenizes it, and dispatches the tokens to their respective commands.
type Interpreter struct {
	api api.API
}

func NewCLI(api api.API) *Interpreter {
	return &Interpreter{
		api: api,
	}
}

// Exec is the new Execute.
func (i *Interpreter) Exec(line string) error {

	// Create the token scanner.
	reader := strings.NewReader(line)
	scanner := lexer.NewScanner(reader)

	// Read the verb of the function. Comments and whitespace are ignored, all
	// tokens other than identifiers throw errors.
	tok, verb := scanner.ScanIgnoreWhitespace()
	switch tok {
	case lexer.TokenEnd, lexer.TokenComment:
		return nil
	case lexer.TokenIdentifier:
		break
	default:
		return fmt.Errorf("Unexpected '%s', expected verb", verb)
	}

	// Instantiate the command.
	cmd := commands.New(verb, i.api)
	if cmd == nil {
		return fmt.Errorf("Not a command: %s", verb)
	}

	// Parse the command into an AST.
	cmd.SetScanner(scanner)
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// Execute the AST.
	return cmd.Exec()
}

// Execute sends scanned tokens to Command instances.
// FIXME: this function is deprecated and must be remove when all Command
// classes have been ported.
func (i *Interpreter) Execute(line string) error {
	var cmd commands.Command
	var err error

	err = i.Exec(line)
	if err != nil {
		return err
	}

	reader := strings.NewReader(line)
	scanner := lexer.NewScanner(reader)

	for {
		class, token := scanner.Scan()

		// First identifier; try to find a command handler
		if cmd == nil {
			switch class {
			case lexer.TokenIdentifier:
				if ctor, ok := commands.Verbs[token]; ok {
					cmd = ctor(i.api)
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
