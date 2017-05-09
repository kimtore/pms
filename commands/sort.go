package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Sort sorts songlists.
type Sort struct {
	api      api.API
	fields   []string
	finished bool
}

func NewSort(api api.API) Command {
	sort := api.Options().StringValue("sort")
	return &Sort{
		api:    api,
		fields: strings.Split(sort, ","),
	}
}

func (cmd *Sort) Execute(t lexer.Token) error {
	var err error

	s := t.String()
	songlistWidget := cmd.api.SonglistWidget()

	switch t.Class {

	case lexer.TokenIdentifier:
		if cmd.finished {
			return fmt.Errorf("Unknown input '%s', expected END", s)
		}
		cmd.fields = strings.Split(s, ",")
		cmd.finished = true

	case lexer.TokenEnd:
		song := songlistWidget.CursorSong()
		err = songlistWidget.Songlist().Sort(cmd.fields)
		songlistWidget.CursorToSong(song)

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}
