package commands

import (
	"fmt"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/widgets"
)

// Add adds songs to MPD's queue.
type Add struct {
	api      api.API
	song     *song.Song
	songlist songlist.Songlist
}

func NewAdd(api api.API) Command {
	return &Add{
		api: api,
	}
}

func (cmd *Add) Execute(class int, s string) error {
	var err error

	switch class {
	case lexer.TokenIdentifier:
		if cmd.song != nil {
			return fmt.Errorf("Cannot add multiple paths on the same command line.")
		}
		cmd.song = song.New()
		cmd.song.SetTags(mpd.Attrs{
			"file": s,
		})

	case lexer.TokenEnd:
		list := cmd.api.SonglistWidget().Songlist()
		queue := cmd.api.Queue()

		switch {
		case cmd.song == nil:
			selection := list.Selection()
			if selection.Len() == 0 {
				return fmt.Errorf("No selection, cannot add without any parameters.")
			}
			err = queue.AddList(selection)
			if err != nil {
				break
			}
			list.ClearSelection()
			cmd.api.Multibar().SetMode(widgets.MultibarModeNormal) // FIXME: remove
			list.MoveCursor(1)
			len := selection.Len()
			if len == 1 {
				song := selection.Songs()[0]
				cmd.api.Message("Added to queue: %s", song.StringTags["file"])
			} else {
				cmd.api.Message("Added %d songs to queue.", len)
			}

		default:
			err = queue.Add(cmd.song)
			if err == nil {
				cmd.api.Message("Added to queue: %s", cmd.song.StringTags["file"])
			}
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}
