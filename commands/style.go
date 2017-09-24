package commands

import (
	"fmt"
	"sort"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/term"
)

// Style manipulates the style table, allowing to set colors and attributes for UI elements.
type Style struct {
	newcommand
	api api.API

	styleKey   string
	styleValue term.Style

	background bool
	foreground bool
}

// NewStyle returns Style.
func NewStyle(api api.API) Command {
	return &Style{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Style) Parse() error {

	// Scan the style key. All names are accepted, even names that are not
	// implemented anywhere.
	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteNames(lit)
	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("Unexpected '%v', expected identifier", lit)
	}
	cmd.styleKey = lit

	// Scan each style attribute.
	for {
		tok, lit := cmd.Scan()

		switch tok {
		case lexer.TokenWhitespace:
			cmd.setTabCompleteStyles("")
			continue
		case lexer.TokenIdentifier:
			break
		case lexer.TokenEnd:
			return nil
		default:
			return fmt.Errorf("Unexpected '%v', expected identifier", lit)
		}

		cmd.setTabCompleteStyles(lit)
		err := cmd.mergeStyle(lit)
		if err != nil {
			return err
		}
	}
}

// Exec implements Command.
func (cmd *Style) Exec() error {
	styleMap := cmd.api.Styles()
	styleMap[cmd.styleKey] = cmd.styleValue
	return nil
}

// setTabCompleteNames sets the tab complete list to the list of available style keys.
func (cmd *Style) setTabCompleteNames(lit string) {
	styleMap := cmd.api.Styles()
	list := make(sort.StringSlice, len(styleMap))
	i := 0
	for key := range styleMap {
		list[i] = key
		i++
	}
	list.Sort()
	cmd.setTabComplete(lit, list)
}

// setTabCompleteStyles sets the tab complete list to available styles.
func (cmd *Style) setTabCompleteStyles(lit string) {
	list := []string{
		"bold",
		"reverse",
		"underline",
	}
	cmd.setTabComplete(lit, list)
}

func (cmd *Style) mergeStyle(lit string) error {
	switch lit {
	case "bold":
		cmd.styleValue = cmd.styleValue.Bold(true)
	case "reverse":
		cmd.styleValue = cmd.styleValue.Reverse(true)
	case "underline":
		cmd.styleValue = cmd.styleValue.Underline(true)
	default:
		if lit[0] == '@' {
			lit = "#" + lit[1:]
		}
		color := term.GetColor(lit)
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

	return nil
}
