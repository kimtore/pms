package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/input/lexer"
)

// Print displays information about the selected song's tags.
type Print struct {
	api  API
	tags []string
}

func NewPrint(api API) Command {
	return &Print{
		api:  api,
		tags: make([]string, 0),
	}
}

func (cmd *Print) Execute(t lexer.Token) error {
	var err error

	s := t.String()

	switch t.Class {
	case lexer.TokenIdentifier:
		if len(cmd.tags) > 0 {
			return fmt.Errorf("Unexpected '%s', expected END")
		}
		cmd.tags = strings.Split(strings.ToLower(s), ",")

	case lexer.TokenEnd:
		if len(cmd.tags) == 0 {
			return fmt.Errorf("Unexpected END, expected list of tags to print")
		}
		songlistWidget := cmd.api.SonglistWidget()
		selection := songlistWidget.Selection()
		switch selection.Len() {
		case 0:
			return fmt.Errorf("Cannot print song tags; no song selected")
		case 1:
			song := selection.Song(0)
			parts := make([]string, 0)
			for _, tag := range cmd.tags {
				msg := ""
				value, ok := song.StringTags[tag]
				if ok {
					msg = fmt.Sprintf("%s: '%s'", tag, value)
				} else {
					msg = fmt.Sprintf("%s: <MISSING>", tag)
				}
				parts = append(parts, msg)
			}
			msg := strings.Join(parts, ", ")
			cmd.api.Message(msg)

		default:
			return fmt.Errorf("Multiple songs selected; cannot print song tags")
		}

		songlistWidget.ClearSelection()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return err
}
