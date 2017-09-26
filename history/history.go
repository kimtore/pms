// package history provides a backlog of text lines that can be navigated.
package history

// History represents a text history that can be navigated through.
type History struct {
	items   []string
	current string
	index   int
}

// New returns History.
func New() *History {
	return &History{
		items: make([]string, 0),
	}
}

// Add adds to the input history.
func (h *History) Add(s string) {
	if len(s) > 0 {
		hl := len(h.items)
		if hl == 0 || h.items[hl-1] != s {
			h.items = append(h.items, s)
		}
	}
	h.Reset(s)
}

// Reset resets the cursor offset to the last position.
func (h *History) Reset(s string) {
	h.index = len(h.items)
	h.current = s
}

// Current returns the current history item.
func (h *History) Current() string {
	if len(h.items) == 0 || h.index >= len(h.items) {
		h.index = len(h.items)
		return h.current
	}
	h.validateIndex()
	return h.items[h.index]
}

// Navigate navigates the history and returns that history item.
func (h *History) Navigate(offset int) string {
	h.index += offset
	return h.Current()
}

// validateIndex ensures that the item index stays within the valid range.
func (h *History) validateIndex() {
	if h.index >= len(h.items) {
		h.index = len(h.items) - 1
	}
	if h.index < 0 {
		h.index = 0
	}
}
