package commands

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/keysequence"
)

// Unbind unmaps a key sequence.
type Unbind struct {
	command
	api api.API
	seq keysequence.KeySequence
}

// NewUnbind returns Unbind.
func NewUnbind(api api.API) Command {
	return &Unbind{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Unbind) Parse() error {

	// Use the key sequence parser for parsing the next token.
	parser := keysequence.NewParser(cmd.S)

	// Parse a valid key sequence from the scanner.
	seq, err := parser.ParseKeySequence()
	if err != nil {
		return err
	}
	cmd.seq = seq

	// Reject any further input
	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Unbind) Exec() error {
	sequencer := cmd.api.Sequencer()
	return sequencer.RemoveBind(cmd.seq)
}
