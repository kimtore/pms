// Package commands contains all functionality that is triggered by the user,
// either through keyboard bindings or the command-line interface. New commands
// such as 'sort', 'add', etc. must be implemented here.
package commands

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/parser"
)

// Verbs contain mappings from strings to Command constructors.
// Make sure to add commands here when implementing them, or they will not be recognized.
var Verbs = map[string]func(api.API) Command{
	"add":       NewAdd,
	"bind":      NewBind,
	"cursor":    NewCursor,
	"inputmode": NewInputMode,
	"isolate":   NewIsolate,
	"list":      NewList,
	"next":      NewNext,
	"pause":     NewPause,
	"play":      NewPlay,
	"prev":      NewPrevious,
	"previous":  NewPrevious,
	"print":     NewPrint,
	"q":         NewQuit,
	"quit":      NewQuit,
	"redraw":    NewRedraw,
	"remove":    NewRemove,
	"se":        NewSet,
	"select":    NewSelect,
	"set":       NewSet,
	"sort":      NewSort,
	"stop":      NewStop,
	"style":     NewStyle,
	"volume":    NewVolume,
}

// Command must be implemented by all commands.
type Command interface {
	// Execute parses the next input token.
	// FIXME: Execute is deprecated
	Execute(class int, s string) error

	// Parse and make an abstract syntax tree. This function MUST NOT have any side effects.
	Parse(*lexer.Scanner) error

	// TabComplete returns a set of tokens that could possibly be used as the next
	// command parameter.
	TabComplete() []string

	// Scanned returns a slice of tokens that have been scanned using Parse().
	Scanned() []parser.Token
}

// command is a helper base class that all commands may use.
type command struct {
	cmdline     string
	tabComplete []string
}

// New returns the Command associated with the given verb.
func New(verb string, a api.API) Command {
	ctor := Verbs[verb]
	if ctor == nil {
		return nil
	}
	return ctor(a)
}

// setTabComplete defines a string slice that will be used for tab completion
// at the current point in parsing.
func (c *command) setTabComplete(s []string) {
	c.tabComplete = s
}

// setTabCompleteEmpty removes all tab completions.
func (c *command) setTabCompleteEmpty() {
	c.setTabComplete([]string{})
}

// TabComplete implements Command.TabComplete.
func (c *command) TabComplete() []string {
	if c.tabComplete == nil {
		// FIXME
		return make([]string, 0)
	}
	return c.tabComplete
}

// Parse implements Command.Parse.
func (c *command) Parse(*lexer.Scanner) error {
	// FIXME
	return nil
}

// Scanned implements Command.Scanned.
func (c *command) Scanned() []parser.Token {
	// FIXME
	return make([]parser.Token, 0)
}
