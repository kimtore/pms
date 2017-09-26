// package term provides terminal emulator events such as keyboard, mouse, and resize.
package term

import (
	"github.com/ambientsound/pms/console"
	termbox "github.com/nsf/termbox-go"
)

// EventType represents an event type such as keyboard, mouse, or resize.
type EventType uint8

// Event type constants.
const (
	EventNone   EventType = iota // no-op or unsupported events.
	EventKey                     // pass-through keypress in direct mode.
	EventMouse                   // any mouse press.
	EventResize                  // terminal resize.
)

// Event represents a terminal event.
type Event struct {
	Type EventType
	Key  KeyPress
}

// Sampler reads all events from the terminal, including keyboard, mouse, and
// resizes. The events are buffered if requested, and sent asynchronously on an
// event channel.
type Sampler struct {
	Events chan Event // event channel, emitting terminal events as they arrive.
}

// NewSampler returns Sampler.
func NewSampler() *Sampler {
	return &Sampler{
		Events: make(chan Event, 1024),
	}
}

// Loop samples keyboard events from the terminal library, and passes them on to the main thread.
func (s *Sampler) Loop() {
	for {
		te := termbox.PollEvent()
		e := s.SampleEvent(te)
		s.Events <- e
	}
}

// sampleEvent samples one event from the terminal library, and returns an Event based on the new internal state.
func (s *Sampler) SampleEvent(te termbox.Event) Event {
	switch te.Type {
	case termbox.EventKey:
		return Event{
			Type: EventKey,
			Key:  ParseKey(te),
		}
	case termbox.EventResize:
		return Event{
			Type: EventResize,
		}
	default:
		console.Log("Ignoring terminal event: %+v", te)
		return Event{}
	}
}
