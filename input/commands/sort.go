package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/widgets"
)

// Sort sorts songlists.
type Sort struct {
	songlistWidget func() *widgets.SonglistWidget
	options        *options.Options
	fields         []string
	finished       bool
}

func NewSort(songlistWidget func() *widgets.SonglistWidget, options *options.Options) *Sort {
	return &Sort{songlistWidget: songlistWidget, options: options}
}

func (cmd *Sort) Reset() {
	sort := cmd.options.StringValue("sort")
	cmd.fields = strings.Split(sort, ",")
	cmd.finished = false
}

func (cmd *Sort) Execute(t lexer.Token) error {
	var err error

	s := t.String()
	songlistWidget := cmd.songlistWidget()

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
