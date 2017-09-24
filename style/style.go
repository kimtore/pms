package style

import (
	"github.com/ambientsound/pms/term"
)

type Stylesheet map[string]term.Style

type Stylable interface {
	Style(string) term.Style
	SetStylesheet(Stylesheet)
	Stylesheet() Stylesheet
}

// Styled implements Stylable
type Styled struct {
	stylesheet Stylesheet
}

func (w *Styled) Style(s string) term.Style {
	return w.stylesheet[s]
}

func (w *Styled) SetStylesheet(stylesheet Stylesheet) {
	w.stylesheet = stylesheet
}

func (w *Styled) Stylesheet() Stylesheet {
	return w.stylesheet
}
