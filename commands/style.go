package commands

import (
	"fmt"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/gdamore/tcell"
)

// Style manipulates the style table, allowing to set colors and attributes for UI elements.
type Style struct {
	command
	api api.API

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

func NewStyle(api api.API) Command {
	return &Style{
		api: api,
	}
}

func (cmd *Style) Execute(class int, s string) error {
	var err error

	switch class {

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
		styleMap := cmd.api.Styles()
		styleMap[cmd.styleKey] = cmd.styleValue

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}
