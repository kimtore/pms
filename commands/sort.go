package commands

import (
	"fmt"
	"github.com/ambientsound/pms/log"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/list"
)

// Sort sorts songlists.
type Sort struct {
	newcommand
	api  api.API
	tags []string
	list list.List
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

	cmd.list = cmd.api.List()
	possibleTags := cmd.list.ColumnNames()

	for {
		tok, lit := cmd.Scan()
		switch tok {
		case lexer.TokenWhitespace:
			// Initialize tab completion
			cmd.setTabComplete("", possibleTags)
			continue

		case lexer.TokenIdentifier:
			// Sort by tags specified on the command line
			cmd.Unscan()
			cmd.tags, err = cmd.ParseTags(possibleTags)
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
	// FIXME: retain cursor
	log.Debugf("sorting list by tags %v", cmd.tags)
	return cmd.list.Sort(cmd.tags)
}
