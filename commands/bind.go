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
	api      api.API
	sentence []string
	token    *parser.KeySequenceToken
}

func NewBind(api api.API) Command {
	return &Bind{
		api:      api,
		sentence: make([]string, 0),
	}
}

func (cmd *Bind) Execute(t lexer.Token) error {
	s := t.String()

	switch t.Class {
	case lexer.TokenIdentifier:
		if cmd.token == nil {
			cmd.token = &parser.KeySequenceToken{}
			err := cmd.token.Parse(t.Runes)
			if err != nil {
				return err
			}
		} else {
			cmd.sentence = append(cmd.sentence, s)
		}

	case lexer.TokenEnd:
		switch {
		case cmd.token == nil:
			return fmt.Errorf("Unexpected END, expected key sequence")
		case len(cmd.sentence) == 0:
			return fmt.Errorf("Unexpected END, expected verb")
		default:
			return cmd.bind()
		}

	default:
		if t.Class != lexer.TokenIdentifier {
			return fmt.Errorf("Unknown input '%s', expected identifier", s)
		}
	}

	return nil
}

func (cmd *Bind) bind() error {
	sentence := strings.Join(cmd.sentence, " ")
	sequencer := cmd.api.Sequencer()
	return sequencer.AddBind(cmd.token.Sequence, sentence)
}
