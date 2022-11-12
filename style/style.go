package style

import (
	"github.com/gdamore/tcell/v2"
)

type Stylesheet map[string]tcell.Style

type Stylable interface {
	Style(string) tcell.Style
	SetStylesheet(Stylesheet)
	Stylesheet() Stylesheet
}

// Styled implements Stylable
type Styled struct {
	stylesheet Stylesheet
}

func (w *Styled) Style(s string) tcell.Style {
	return w.stylesheet[s]
}

func (w *Styled) SetStylesheet(stylesheet Stylesheet) {
	w.stylesheet = stylesheet
}

func (w *Styled) Stylesheet() Stylesheet {
	return w.stylesheet
}
