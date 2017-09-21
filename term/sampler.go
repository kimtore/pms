// package term provides terminal emulator events such as keyboard, mouse, and resize.
package term

import (
	"github.com/ambientsound/pms/console"
	termbox "github.com/nsf/termbox-go"
)

type EventType uint8

type Mode uint8

// Event types.
const (
	EventNone         EventType = iota // no-op or unsupported events.
	EventKey                           // pass-through keypress in direct mode.
	EventBuffer                        // buffer or cursor is updated.
	EventBufferReturn                  // user presses <Return> during buffered mode.
	EventBufferCancel                  // user cancels buffered mode.
	EventMouse                         // any mouse press.
	EventResize                        // terminal resize.
)

// Input modes.
const (
	ModeDirect   Mode = iota // Direct mode will pass along any keyboard events as they arrive.
	ModeBuffered             // Buffered mode will keep a buffer of text, and provide tab completion, history, and navigation.
)

// KeyPress represents a single keypress.
type KeyPress struct {
	Ch  rune
	Key termbox.Key
	Mod termbox.Modifier
}

// Event represents a
type Event struct {
	Type EventType
	Key  KeyPress
}

// Sampler holds buffered input state.
type Sampler struct {
	Events chan Event
	mode   Mode
}

// New returns Sampler. Input events are retrieved asynchronously using the Events channel.
func New() *Sampler {
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
			Key: KeyPress{
				Ch:  te.Ch,
				Key: te.Key,
				Mod: te.Mod,
			},
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
