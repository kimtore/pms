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
	EventNone         EventType = iota // no-op or unsupported events.
	EventKey                           // pass-through keypress in direct mode.
	EventBuffer                        // buffer or cursor is updated.
	EventBufferReturn                  // user presses <Return> during buffered mode.
	EventBufferCancel                  // user cancels buffered mode.
	EventMouse                         // any mouse press.
	EventResize                        // terminal resize.
)

// Mode represents an input mode, such as direct or buffered.
type Mode uint8

// Input modes.
const (
	ModeDirect   Mode = iota // Direct mode will pass along any keyboard events as they arrive.
	ModeBuffered             // Buffered mode will keep a buffer of text, and provide tab completion, history, and navigation.
)

// Event represents a
type Event struct {
	Type EventType
	Key  KeyPress
}

// Sampler reads all events from the terminal, including keyboard, mouse, and
// resizes. The events are buffered if requested, and sent asynchronously on an
// event channel.
type Sampler struct {
	Events chan Event
	mode   Mode
}

// NewSampler returns Sampler.
func NewSampler() *Sampler {
	return &Sampler{
		Events: make(chan Event, 1024),
	}
}

// bufferKey reads a keypress and returns a suitable event.
func (s *Sampler) bufferKey(te termbox.Event) Event {
	switch s.mode {
	case ModeDirect:
		return Event{
			Type: EventKey,
			Key:  ParseKey(te),
		}
	default:
		return Event{}
	}
}

// resize returns a terminal resize event.
func (s *Sampler) resize() Event {
	return Event{
		Type: EventResize,
	}
}

// sampleEvent samples one event from the terminal library, and returns an Event based on the new internal state.
func (s *Sampler) SampleEvent(te termbox.Event) Event {
	switch te.Type {
	case termbox.EventKey:
		return s.bufferKey(te)
	case termbox.EventResize:
		return s.resize()
	default:
		console.Log("Ignoring terminal event: %+v", te)
		return Event{}
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
