package keys

import (
	"fmt"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/keysequence"
	"github.com/gdamore/tcell/v2"
)

// Binding holds a parsed, user provided key sequence.
type Binding struct {
	Command  string
	Sequence keysequence.KeySequence
}

// Sequencer holds all the keyboard bindings and their action mappings.
type Sequencer struct {
	binds []Binding
	event *tcell.EventKey
	input keysequence.KeySequence
}

// NewSequencer returns Sequencer.
func NewSequencer() *Sequencer {
	return &Sequencer{
		binds: make([]Binding, 0),
		input: make(keysequence.KeySequence, 0),
	}
}

// AddBind creates a new key mapping.
func (s *Sequencer) AddBind(seq keysequence.KeySequence, command string) error {
	if s.dupes(seq) {
		return fmt.Errorf("Can't bind: conflicting with already bound key sequence")
	}
	s.binds = append(s.binds, Binding{Sequence: seq, Command: command})
	return nil
}

// RemoveBind removes a key mapping.
func (s *Sequencer) RemoveBind(seq keysequence.KeySequence) error {
	for i := range s.binds {
		if keysequence.Compare(s.binds[i].Sequence, seq) {
			// Overwrite this position with the last in the list
			s.binds[i] = s.binds[len(s.binds)-1]

			// Truncate to remove the (now duplicate) last entry
			s.binds = s.binds[:len(s.binds)-1]

			return nil
		}
	}

	return fmt.Errorf("Can't unbind: sequence not bound")
}

// KeyInput feeds a keypress to the sequencer. Returns true if there is one match or more, or false if there is no match.
func (s *Sequencer) KeyInput(ev *tcell.EventKey) bool {
	console.Log("Key event: %s", keysequence.FormatKey(ev))
	s.input = append(s.input, ev)
	if len(s.find(s.input)) == 0 {
		s.input = make(keysequence.KeySequence, 0)
		return false
	}
	return true
}

// String returns the current input sequence as a string.
func (s *Sequencer) String() string {
	return keysequence.Format(s.input)
}

// dupes returns true if binding the given key event sequence will conflict with any other bound sequences.
func (s *Sequencer) dupes(seq keysequence.KeySequence) bool {
	matches := s.find(seq)
	return len(matches) > 0
}

// find returns a list of potential matches to key bindings.
func (s *Sequencer) find(seq keysequence.KeySequence) []Binding {
	binds := make([]Binding, 0)
	for i := range s.binds {
		if keysequence.StartsWith(s.binds[i].Sequence, seq) {
			binds = append(binds, s.binds[i])
		}
	}
	return binds
}

// Match returns a key binding if the current input sequence is found.
func (s *Sequencer) Match() *Binding {
	binds := s.find(s.input)
	if len(binds) != 1 {
		return nil
	}
	b := binds[0]
	//console.Log("Possible match found: %+v ||| %+v", b.Sequence, s.input)
	if !keysequence.Compare(b.Sequence, s.input) {
		return nil
	}
	//console.Log("Match found: %+v", b)
	s.input = make(keysequence.KeySequence, 0)
	return &b
}
