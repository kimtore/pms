package commands

import (
	"fmt"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/widgets"
)

// Add adds songs to MPD's queue.
type Add struct {
	songlistWidget func() *widgets.SonglistWidget
	queue          func() *songlist.Queue
	song           *song.Song
	songlist       songlist.Songlist
}

func NewAdd(songlistWidget func() *widgets.SonglistWidget, queue func() *songlist.Queue) *Add {
	return &Add{songlistWidget: songlistWidget, queue: queue}
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
		songlistWidget := cmd.songlistWidget()
		queue := cmd.queue()

		switch {
		case cmd.song == nil:
			selection := songlistWidget.Selection()
			if selection.Len() == 0 {
				return fmt.Errorf("No selection, cannot add without any parameters.")
			}
			err := queue.AddList(selection)
			if err != nil {
				break
			}
			if err == nil {
				songlistWidget.ClearSelection()
				songlistWidget.MoveCursor(1)
			}
		default:
			err = queue.Add(cmd.song)
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return err
}
