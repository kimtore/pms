package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/widgets"

	"github.com/fhs/gompd/mpd"
)

// Play plays songs in the MPD playlist.
type Play struct {
	songlistWidget *widgets.SongListWidget
	mpdClient      func() *mpd.Client
	song           *song.Song
	id             int
	pos            int
}

func NewPlay(songlistWidget *widgets.SongListWidget, mpdClient func() *mpd.Client) *Play {
	return &Play{songlistWidget: songlistWidget, mpdClient: mpdClient}
}

func (cmd *Play) Reset() {
	cmd.song = nil
	cmd.pos = -1
}

func (cmd *Play) Execute(t lexer.Token) error {
	var err error

	s := string(t.Runes)

	switch t.Class {
	case lexer.TokenIdentifier:
		switch s {
		case "cursor":
			cmd.song = cmd.songlistWidget.CursorSong()
			if cmd.song == nil {
				return fmt.Errorf("Cannot play: no song under cursor")
			}
		default:
			return nil
		}

	case lexer.TokenEnd:
		client := cmd.mpdClient()
		if client == nil {
			return fmt.Errorf("Cannot play: not connected to MPD")
		}

		if cmd.song == nil {
			err = client.Play(-1)
			return err
		}

		id, err := client.AddID(cmd.song.TagString("file"), -1)
		if err != nil {
			return err
		}

		err = client.PlayID(id)
		return err

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return nil
}
