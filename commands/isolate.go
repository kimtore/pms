package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Isolate searches for songs that have similar tags as the selection.
type Isolate struct {
	api  api.API
	tags []string
}

func NewIsolate(api api.API) Command {
	return &Isolate{
		api:  api,
		tags: make([]string, 0),
	}
}

func (cmd *Isolate) Execute(class int, s string) error {
	var err error

	switch class {
	case lexer.TokenIdentifier:
		if len(cmd.tags) != 0 {
			return fmt.Errorf("Unexpected '%s', expected END.", s)
		}
		cmd.tags = strings.Split(strings.ToLower(s), ",")

	case lexer.TokenEnd:
		if len(cmd.tags) == 0 {
			return fmt.Errorf("Unexpected END, expected comma-separated tags to isolate by.")
		}

		index := cmd.api.Index()
		if index == nil {
			return fmt.Errorf("Search index is not operational.")
		}

		songlistWidget := cmd.api.SonglistWidget()
		list := cmd.api.Songlist()
		selection := list.Selection()
		song := list.CursorSong()

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

		sort := cmd.api.Options().StringValue("sort")
		fields := strings.Split(sort, ",")
		result.Sort(fields)

		list.ClearSelection()
		songlistWidget.AddSonglist(result)
		songlistWidget.SetSonglist(result)
		list.CursorToSong(song)

	default:
		return fmt.Errorf("Unknown input '%s', expected END.", s)
	}

	return err
}
