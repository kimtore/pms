package commands

import "github.com/ambientsound/pms/input/lexer"

type Command interface {
	// Parse the next input token
	Execute(t lexer.Token) error
}
