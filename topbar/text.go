package topbar

// Text draws a literal text string.
type Text struct {
	fragment
	t string
}

func NewText(s string) Fragment {
	return &Text{
		t: s,
	}
}

func (w *Text) Width() int {
	return len(w.t)
}

func (w *Text) Draw(x, y int) int {
	return w.drawNextString(x, y, w.t, w.Style("topbar"))
}
