package commands

import (
	"fmt"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/input/parser"
)

type runeString []rune

// Bind maps a key sequence to the execution of a command.
type Bind struct {
	sentence []runeString
	token    *parser.KeySequenceToken
}

func NewBind() *Bind {
	p := &Bind{}
	p.Reset()
	return p
}

func (p *Bind) Reset() {
	p.token = nil
	p.sentence = make([]runeString, 0)
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
			p.sentence = append(p.sentence, t.Runes)
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
	console.Log("Binding key input sequence %v => %v", p.token, p.sentence)
	return nil
}
