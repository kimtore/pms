package keys

import (
	"fmt"

	"github.com/ambientsound/pms/input/parser"
)

type Input struct {
	Multiplier int
	Command    string
	Sequence   parser.KeyEvents
}

type Sequencer struct {
	binds []Input
	event parser.KeyEvent
	input Input
}

func NewSequencer() *Sequencer {
	s := &Sequencer{}
	s.binds = make([]Input, 0)
	return s
}

func (s *Sequencer) AddBind(seq parser.KeyEvents, command string) error {
	if s.dupes(seq) {
		return fmt.Errorf("Can't bind '%s': conflicting with already bound key sequence", seq.String())
	}
	s.binds = append(s.binds, Input{Sequence: seq, Command: command})
	return nil
}

// KeyInput feeds a keypress to the sequencer. Returns true if there is one match or more, or false if there is no match.
func (s *Sequencer) KeyInput(ev parser.KeyEvent) bool {
	s.input.Sequence = append(s.input.Sequence, ev)
	if len(s.find(s.input.Sequence)) == 0 {
		s.input = Input{}
		return false
	}
	return true
}

// String returns the current input sequence as a string.
func (s *Sequencer) String() string {
	return s.input.Sequence.String()
}

// dupes returns true if binding the given key event sequence will conflict with any other bound sequences.
func (s *Sequencer) dupes(seq parser.KeyEvents) bool {
	for i := range s.binds {
		if s.binds[i].Sequence.StartsWith(seq) || seq.StartsWith(s.binds[i].Sequence) {
			return true
		}
	}
	return false
}

func (s *Sequencer) find(seq parser.KeyEvents) []Input {
	binds := make([]Input, 0)
	for i := range s.binds {
		if s.binds[i].Sequence.StartsWith(seq) {
			binds = append(binds, s.binds[i])
		}
	}
	return binds
}

func (s *Sequencer) Match() *Input {
	binds := s.find(s.input.Sequence)
	for i := range binds {
		if binds[i].Sequence.Equals(s.input.Sequence) {
			x := &Input{
				Multiplier: s.input.Multiplier,
				Command:    binds[i].Command,
			}
			s.input = Input{}
			return x
		}
	}
	return nil
}
