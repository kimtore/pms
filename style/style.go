package style

import (
	"github.com/ambientsound/pms/log"
	"github.com/gdamore/tcell"
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

func (w *Styled) MessageStyle(msg log.Message) tcell.Style {
	switch msg.Level {
	case log.InfoLevel:
		return w.Style("statusbar")
	case log.ErrorLevel:
		return w.Style("errorText")
	default:
		return w.Style("default")
	}
}
