package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Print displays information about the selected song's tags.
type Print struct {
	command
	api  api.API
	tags []string
}

func NewPrint(api api.API) Command {
	return &Print{
		api:  api,
		tags: make([]string, 0),
	}
}

func (cmd *Print) Execute(class int, s string) error {
	var err error

	switch class {
	case lexer.TokenIdentifier:
		if len(cmd.tags) > 0 {
			return fmt.Errorf("Unexpected '%s', expected END", s)
		}
		cmd.tags = strings.Split(strings.ToLower(s), ",")

	case lexer.TokenEnd:
		if len(cmd.tags) == 0 {
			return fmt.Errorf("Unexpected END, expected list of tags to print")
		}
		list := cmd.api.Songlist()
		selection := list.Selection()
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
					value = strings.ReplaceAll(value, "%", "%%")
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

		list.ClearSelection()

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}
