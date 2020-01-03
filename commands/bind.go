package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/keysequence"
)

// Bind maps a key sequence to the execution of a command.
type Bind struct {
	command
	api      api.API
	sentence string
	context  string
	seq      keysequence.KeySequence
}

// NewBind returns Bind.
func NewBind(api api.API) Command {
	return &Bind{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Bind) Parse() error {

	// Bind keyboard sequence to a specific program context.
	err := cmd.ParseContext()
	if err != nil {
		return err
	}

	// Use the key sequence parser for parsing the next token.
	parser := keysequence.NewParser(cmd.S)

	// Parse a valid key sequence from the scanner.
	seq, err := parser.ParseKeySequence()
	if err != nil {
		return err
	}
	cmd.seq = seq

	// Treat the rest of the line as the literal action to execute when the bind succeeds.
	sentence := make([]string, 0, 32)
	for {
		tok, lit := cmd.Scan()
		if tok == lexer.TokenEnd {
			break
		} else if tok == lexer.TokenIdentifier {
			// Quote identifiers?
		}
		sentence = append(sentence, lit)
	}

	if len(sentence) == 0 {
		return fmt.Errorf("Unexpected END, expected identifier")
	}

	cmd.sentence = strings.Join(sentence, "")
	return nil
}

func (cmd *Bind) ParseContext() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabComplete(lit, contexts)

	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	switch lit {
	case GlobalContext, ListContext, TracklistContext:
		cmd.context = lit
		cmd.setTabCompleteEmpty()
		return nil
	default:
		return fmt.Errorf("unexpected '%s', expected one of %v", lit, contexts)
	}
}

// Exec implements Command.
func (cmd *Bind) Exec() error {
	sequencer := cmd.api.Sequencer()
	return sequencer.AddBind(cmd.seq, cmd.sentence)
}
