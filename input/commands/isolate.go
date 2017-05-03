package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/widgets"
)

// Isolate searches for songs that have similar tags as the selection.
type Isolate struct {
	messages       chan string
	songlistWidget func() *widgets.SonglistWidget
	index          func() *index.Index
	options        *options.Options
	tags           []string
}

func NewIsolate(
	messages chan string,
	songlistWidget func() *widgets.SonglistWidget,
	index func() *index.Index,
	options *options.Options) *Isolate {

	return &Isolate{
		messages:       messages,
		songlistWidget: songlistWidget,
		index:          index,
		options:        options,
	}
}

func (cmd *Isolate) Reset() {
	cmd.tags = make([]string, 0)
}

func (cmd *Isolate) Execute(t lexer.Token) error {
	var err error

	s := t.String()

	switch t.Class {
	case lexer.TokenIdentifier:
		if len(cmd.tags) != 0 {
			return fmt.Errorf("Unexpected '%s', expected END.", s)
		}
		cmd.tags = strings.Split(strings.ToLower(s), ",")

	case lexer.TokenEnd:
		if len(cmd.tags) == 0 {
			return fmt.Errorf("Unexpected END, expected comma-separated tags to isolate by.")
		}

		index := cmd.index()
		if index == nil {
			return fmt.Errorf("Search index is not operational.")
		}

		songlistWidget := cmd.songlistWidget()
		selection := songlistWidget.Selection()
		song := songlistWidget.CursorSong()

		if selection.Len() == 0 {
			return fmt.Errorf("No selection, cannot isolate in empty songlist.")
		}

		result, err := index.Isolate(selection, cmd.tags)
		if err != nil {
			return err
		}

		if result.Len() == 0 {
			return fmt.Errorf("No results found when isolating by '%s'.", strings.Join(cmd.tags, ","))
		}

		sort := cmd.options.StringValue("sort")
		fields := strings.Split(sort, ",")
		err = result.Sort(fields)

		songlistWidget.ClearSelection()
		songlistWidget.AddSonglist(result)
		songlistWidget.SetSonglist(result)
		songlistWidget.CursorToSong(song)

	default:
		return fmt.Errorf("Unknown input '%s', expected END.", s)
	}

	return err
}
