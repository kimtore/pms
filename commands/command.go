// Package commands contains all functionality that is triggered by the user,
// either through keyboard bindings or the command-line interface. New commands
// such as 'sort', 'add', etc. must be implemented here.
package commands

import "github.com/ambientsound/pms/input/lexer"

// Command must be implemented by all commands.
type Command interface {
	// Parse the next input token
	Execute(class int, s string) error

	// Parse and make an abstract syntax tree
	Parse(*lexer.Scanner) error
}

type command struct{}

// Parse implements Command.Parse.
func (c *command) Parse(*lexer.Scanner) error {
	// FIXME
	return nil
}
