package keys

import (
	"fmt"
	"github.com/ambientsound/pms/log"

	"github.com/ambientsound/pms/keysequence"
	"github.com/gdamore/tcell"
)

// Binding holds a parsed, user provided key sequence.
type Binding struct {
	Command  string
	Context  string
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
func (s *Sequencer) AddBind(binding Binding) error {
	if s.dupes(binding) {
		return fmt.Errorf("can't bind: conflicting with already bound key sequence")
	}
	s.binds = append(s.binds, binding)
	return nil
}

// RemoveBind removes a key mapping.
func (s *Sequencer) RemoveBind(context string, seq keysequence.KeySequence) error {
	for i := range s.binds {
		if context != s.binds[i].Context {
			continue
		}
		if !keysequence.Compare(s.binds[i].Sequence, seq) {
			continue
		}

		// Overwrite this position with the last in the list
		s.binds[i] = s.binds[len(s.binds)-1]

		// Truncate to remove the (now duplicate) last entry
		s.binds = s.binds[:len(s.binds)-1]

		return nil
	}

	return fmt.Errorf("can't unbind: sequence not bound")
}

// KeyInput feeds a keypress to the sequencer. Returns true if there is one match or more, or false if there is no match.
func (s *Sequencer) KeyInput(ev *tcell.EventKey, contexts []string) bool {
	log.Debugf("Key event: %s", keysequence.FormatKey(ev))
	s.input = append(s.input, ev)
	if len(s.findAll(s.input, contexts)) == 0 {
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
func (s *Sequencer) dupes(bind Binding) bool {
	matches := s.find(bind.Sequence, bind.Context)
	return len(matches) > 0
}

func (s *Sequencer) findAll(seq keysequence.KeySequence, contexts []string) []Binding {
	binds := make([]Binding, 0)
	for _, context := range contexts {
		binds = append(binds, s.find(seq, context)...)
	}
	return binds
}

// find returns a list of potential matches to key bindings.
func (s *Sequencer) find(seq keysequence.KeySequence, context string) []Binding {
	binds := make([]Binding, 0)
	for i := range s.binds {
		if s.binds[i].Context == context && keysequence.StartsWith(s.binds[i].Sequence, seq) {
			binds = append(binds, s.binds[i])
		}
	}
	return binds
}

// Match returns a key binding if the current input sequence is found.
func (s *Sequencer) Match(contexts []string) *Binding {
	binds := s.findAll(s.input, contexts)
	if len(binds) != 1 {
		return nil
	}
	b := binds[0]
	if !keysequence.Compare(b.Sequence, s.input) {
		return nil
	}
	s.input = make(keysequence.KeySequence, 0)
	return &b
}
