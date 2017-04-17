package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/input/parser"
)

// Bind maps a key sequence to the execution of a command.
type Bind struct {
	sentence  []string
	token     *parser.KeySequenceToken
	sequencer *keys.Sequencer
}

func NewBind(s *keys.Sequencer) *Bind {
	p := &Bind{sequencer: s}
	p.Reset()
	return p
}

func (p *Bind) Reset() {
	p.token = nil
	p.sentence = make([]string, 0)
}

func (p *Bind) Execute(t lexer.Token) error {

	switch t.Class {
	case lexer.TokenIdentifier:
		if p.token == nil {
			p.token = &parser.KeySequenceToken{}
			err := p.token.Parse(t.Runes)
			if err != nil {
				return err
			}
		} else {
			p.sentence = append(p.sentence, string(t.Runes))
		}

	case lexer.TokenEnd:
		switch {
		case p.token == nil:
			return fmt.Errorf("Unexpected END, expected key sequence")
		case len(p.sentence) == 0:
			return fmt.Errorf("Unexpected END, expected verb")
		default:
			return p.bind()
		}

	default:
		if t.Class != lexer.TokenIdentifier {
			return fmt.Errorf("Unknown input '%s', expected identifier", string(t.Runes))
		}
	}

	return nil
}

func (p *Bind) bind() error {
	sentence := strings.Join(p.sentence, " ")
	return p.sequencer.AddBind(p.token.Sequence, sentence)
}
