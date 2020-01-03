package commands

import (
	"fmt"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/log"
	"strings"

	"github.com/ambientsound/pms/api"
)

// Print displays information about the selected song's tags.
type Print struct {
	command
	api  api.API
	row  list.Row
	tags []string
}

func NewPrint(api api.API) Command {
	return &Print{
		api:  api,
		tags: make([]string, 0),
	}
}

func (cmd *Print) Parse() error {
	var err error

	lst := cmd.api.List()
	if lst.Len() == 0 {
		return fmt.Errorf("cannot print anything for empty lists")
	}

	cmd.row = lst.Row(lst.Cursor())

	tok, lit := cmd.ScanIgnoreWhitespace()
	if tok == lexer.TokenEnd {
		cmd.tags = cmd.row.Keys()
		cmd.setTabComplete(lit, cmd.tags)
	} else {
		cmd.Unscan()
		cmd.tags, err = cmd.ParseTags(cmd.row.Keys())
	}

	return err
}

func (cmd *Print) Exec() error {
	parts := make([]string, 0)

	for _, tag := range cmd.tags {
		msg := ""
		value, ok := cmd.row[tag]
		if ok {
			msg = fmt.Sprintf("%s: '%s'", tag, value)
		} else {
			msg = fmt.Sprintf("%s: <MISSING>", tag)
		}
		parts = append(parts, msg)
	}

	log.Infof(strings.Join(parts, ", "))

	return nil
}
