package commands

import (
	"fmt"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
)

// Add adds songs to MPD's queue.
type Add struct {
	Base
	song     *song.Song
	songlist songlist.Songlist
}

func NewAdd(b Base) Command {
	return &Add{
		Base: b,
	}
}

func (cmd *Add) Reset() {
	cmd.song = nil
	cmd.songlist = nil
}

func (cmd *Add) Execute(t lexer.Token) error {
	var err error

	switch t.Class {
	case lexer.TokenIdentifier:
		if cmd.song != nil {
			return fmt.Errorf("Cannot add multiple paths on the same command line.")
		}
		cmd.song = song.New()
		cmd.song.SetTags(mpd.Attrs{
			"file": t.String(),
		})

	case lexer.TokenEnd:
		songlistWidget := cmd.SonglistWidget()
		queue := cmd.CurrentQueue()

		switch {
		case cmd.song == nil:
			selection := songlistWidget.Selection()
			if selection.Len() == 0 {
				return fmt.Errorf("No selection, cannot add without any parameters.")
			}
			err = queue.AddList(selection)
			if err != nil {
				break
			}
			songlistWidget.ClearSelection()
			songlistWidget.MoveCursor(1)
			len := selection.Len()
			if len == 1 {
				song := selection.Songs()[0]
				cmd.EventMessage <- message.Format("Added to queue: %s", song.StringTags["file"])
			} else {
				cmd.EventMessage <- message.Format("Added %d songs to queue.", len)
			}

		default:
			err = queue.Add(cmd.song)
			if err == nil {
				cmd.EventMessage <- message.Format("Added to queue: %s", cmd.song.StringTags["file"])
			}
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return err
}
