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
		return fmt.Errorf("Can't bind: conflicting with already bound key sequence")
	}
	s.binds = append(s.binds, Input{Sequence: seq, Command: command})
	return nil
}

func (s *Sequencer) KeyInput(ev parser.KeyEvent) error {
	s.input.Sequence = append(s.input.Sequence, ev)
	if len(s.find(s.input.Sequence)) == 0 {
		s.input = Input{}
		return fmt.Errorf("Key sequence not bound")
	}
	return nil
}

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
