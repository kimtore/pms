package commands

import (
	"fmt"
	"github.com/ambientsound/pms/input/keys"
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
	var err error

	// Bind keyboard sequence to a specific program context.
	cmd.context, err = cmd.ParseContext()
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

// Exec implements Command.
func (cmd *Bind) Exec() error {
	sequencer := cmd.api.Sequencer()
	return sequencer.AddBind(keys.Binding{
		Context:  cmd.context,
		Sequence: cmd.seq,
		Command:  cmd.sentence,
	})
}
