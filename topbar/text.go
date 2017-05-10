package topbar

// Text draws a literal text string.
type Text struct {
	text string
}

func NewText(s string) Fragment {
	return &Text{text: s}
}

func (w *Text) Text() (string, string) {
	return w.text, `topbar`
}
