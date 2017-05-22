package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Sort sorts songlists.
type Sort struct {
	newcommand
	api  api.API
	tags []string
}

// NewSort returns Sort.
func NewSort(api api.API) Command {
	return &Sort{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Sort) Parse() error {
	var err error

	// For tab completion
	list := cmd.api.Songlist()
	song := list.CursorSong()

	for {
		tok, lit := cmd.Scan()
		switch tok {
		case lexer.TokenWhitespace:
			// Initialize tab completion
			cmd.setTabCompleteTag("", song)
			continue

		case lexer.TokenIdentifier:
			// Sort by tags specified on the command line
			cmd.Unscan()
			cmd.tags, err = cmd.ParseTags(song)
			return err

		case lexer.TokenEnd:
			// Sort by default tags
			sort := cmd.api.Options().StringValue("sort")
			cmd.tags = strings.Split(sort, ",")
			return nil

		default:
			return fmt.Errorf("Unexpected %v, expected tag", lit)
		}
	}
}

// Exec implements Command.
func (cmd *Sort) Exec() error {
	list := cmd.api.Songlist()
	song := list.CursorSong()
	err := list.Sort(cmd.tags)
	list.CursorToSong(song)
	return err
}
