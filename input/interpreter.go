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

// Exec scans an input line, finds the verb in the command directory,
// and hands execution over to the command.
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
		return fmt.Errorf("unexpected '%s', expected verb", verb)
	}

	// Instantiate the command.
	cmd := commands.New(verb, i.api)
	if cmd == nil {
		return fmt.Errorf("not a command: %s", verb)
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
