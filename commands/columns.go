package commands

import (
	"github.com/ambientsound/pms/api"
)

// Columns sets which column headers should be visible for the current list.
type Columns struct {
	newcommand
	api  api.API
	tags []string
}

// NewColumns returns Columns.
func NewColumns(api api.API) Command {
	return &Columns{
		api:  api,
		tags: make([]string, 0),
	}
}

// Parse implements Command.
func (cmd *Columns) Parse() error {
	var err error
	list := cmd.api.List()
	cmd.tags, err = cmd.ParseTags(list.ColumnNames())
	return err
}

// Exec implements Command.
func (cmd *Columns) Exec() error {
	cmd.api.UI().TableWidget().SetColumns(cmd.tags)
	return nil
}
