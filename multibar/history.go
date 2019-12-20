package multibar

import (
	"github.com/ambientsound/pms/log"
)

// history represents a text history that can be navigated through.
type history struct {
	items   []string
	current string
	index   int
}

func NewHistory() *history {
	return &history{
		items: make([]string, 0),
	}
}

// Add adds to the input history.
func (h *history) Add(s string) {
	if len(s) > 0 {
		hl := len(h.items)
		if hl == 0 || h.items[hl-1] != s {
			h.items = append(h.items, s)
		}
	}
	h.Reset(s)
}

// Reset resets the cursor offset to the last position.
func (h *history) Reset(s string) {
	h.index = len(h.items)
	h.current = s
}

// Current returns the current history item.
func (h *history) Current() string {
	if len(h.items) == 0 || h.index >= len(h.items) {
		log.Debugf("multibar: want index %d, returning current string '%s'", h.index, h.current)
		h.index = len(h.items)
		return h.current
	}
	h.validateIndex()
	log.Debugf("multibar: history returning index %d", h.index)
	return h.items[h.index]
}

// Navigate navigates the history and returns that history item.
func (h *history) Navigate(offset int) string {
	h.index += offset
	return h.Current()
}

// validateIndex ensures that the item index stays within the valid range.
func (h *history) validateIndex() {
	if h.index >= len(h.items) {
		h.index = len(h.items) - 1
	}
	if h.index < 0 {
		h.index = 0
	}
}
