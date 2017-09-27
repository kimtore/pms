package style

type Stylable interface {
	Style(string) Style
	SetStylesheet(Stylesheet)
	Stylesheet() Stylesheet
}

// Styled implements Stylable
type Styled struct {
	stylesheet Stylesheet
}

func (w *Styled) Style(s string) Style {
	return w.stylesheet[s]
}

func (w *Styled) SetStylesheet(stylesheet Stylesheet) {
	w.stylesheet = stylesheet
}

func (w *Styled) Stylesheet() Stylesheet {
	return w.stylesheet
}
