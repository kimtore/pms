package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
)

// Quit exits the program.
type Quit struct {
	signal chan int
}

func NewQuit(signal chan int) *Quit {
	return &Quit{signal: signal}
}

func (cmd *Quit) Reset() {
}

func (cmd *Quit) Execute(t lexer.Token) error {
	switch t.Class {
	case lexer.TokenEnd:
		cmd.signal <- 0
		return nil
	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}
}
