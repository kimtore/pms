package topbar

// Text draws a literal text string.
type Text struct {
	text string
}

// NewText returns Text.
func NewText(s string) Fragment {
	return &Text{text: s}
}

// Text implements Fragment.
func (w *Text) Text() (string, string) {
	return w.text, `topbar`
}
