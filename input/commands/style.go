package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/widgets"
	"github.com/gdamore/tcell"
)

// Style manipulates the style table, allowing to set colors and attributes for UI elements.
type Style struct {
	styleMap   widgets.StyleMap
	styleKey   string
	styleValue tcell.Style
	styled     bool

	background bool
	blink      bool
	bold       bool
	dim        bool
	foreground bool
	reverse    bool
	underline  bool
}

func NewStyle(styleMap widgets.StyleMap) *Style {
	return &Style{
		styleMap: styleMap,
	}
}

func (cmd *Style) Reset() {
	cmd.styleKey = ""
	cmd.styleValue = tcell.StyleDefault
	cmd.styled = false

	cmd.background = false
	cmd.blink = false
	cmd.bold = false
	cmd.dim = false
	cmd.foreground = false
	cmd.reverse = false
	cmd.underline = false
}

func (cmd *Style) Execute(t lexer.Token) error {
	var err error

	s := t.String()

	switch t.Class {

	case lexer.TokenIdentifier:
		if len(cmd.styleKey) == 0 {
			cmd.styleKey = s
			return nil
		}

		switch s {
		case "blink":
			cmd.styleValue = cmd.styleValue.Blink(true)
		case "bold":
			cmd.styleValue = cmd.styleValue.Bold(true)
		case "dim":
			cmd.styleValue = cmd.styleValue.Dim(true)
		case "reverse":
			cmd.styleValue = cmd.styleValue.Reverse(true)
		case "underline":
			cmd.styleValue = cmd.styleValue.Underline(true)
		default:
			if s[0] == '@' {
				s = "#" + s[1:]
			}
			color := tcell.GetColor(s)
			switch {
			case !cmd.foreground:
				cmd.styleValue = cmd.styleValue.Foreground(color)
				cmd.foreground = true
			case !cmd.background:
				cmd.styleValue = cmd.styleValue.Background(color)
				cmd.background = true
			default:
				return fmt.Errorf("Only two color values are allowed per style.")
			}
		}

		cmd.styled = true

	case lexer.TokenEnd:
		if !cmd.styled {
			return fmt.Errorf("Unexpected END, expected style attribute")
		}
		cmd.styleMap[cmd.styleKey] = cmd.styleValue

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}
