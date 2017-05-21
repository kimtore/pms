package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/input/parser"
)

// Bind maps a key sequence to the execution of a command.
type Bind struct {
	newcommand
	api      api.API
	sentence string
	token    *parser.KeySequenceToken
}

// NewBind returns Bind.
func NewBind(api api.API) Command {
	return &Bind{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Bind) Parse() error {

	// Scan the first key, which might be whitespace.
	tok, lit := cmd.ScanIgnoreWhitespace()

	// Scan all non-whitespace and non-end tokens into a string sequence.
	seq := ""
Sequence:
	for {
		switch tok {
		case lexer.TokenWhitespace:
			break Sequence
		case lexer.TokenEnd:
			if len(seq) == 0 {
				return fmt.Errorf("Unexpected END, expected key sequence")
			}
			return fmt.Errorf("Unexpected END, expected verb")
		}
		seq += lit
		tok, lit = cmd.Scan()
	}

	// Parse the key sequence.
	cmd.token = &parser.KeySequenceToken{}
	err := cmd.token.Parse([]rune(seq))
	if err != nil {
		return err
	}

	// Treat the rest of the line as the literal action to execute when the bind succeeds.
	sentence := make([]string, 0, 32)
	for {
		tok, lit = cmd.Scan()
		if tok == lexer.TokenEnd {
			break
		} else if tok == lexer.TokenIdentifier {
			// Quote identifiers?
		}
		sentence = append(sentence, lit)
	}
	cmd.sentence = strings.Join(sentence, "")

	// Accept no more input at this point.
	return nil
}

// Exec implements Command.
func (cmd *Bind) Exec() error {
	sequencer := cmd.api.Sequencer()
	return sequencer.AddBind(cmd.token.Sequence, cmd.sentence)
}
