package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/widgets"

	"github.com/ambientsound/gompd/mpd"
)

// Add adds songs to MPD's queue.
type Add struct {
	songlistWidget func() *widgets.SonglistWidget
	mpdClient      func() *mpd.Client
	song           *song.Song
}

func NewAdd(songlistWidget func() *widgets.SonglistWidget, mpdClient func() *mpd.Client) *Add {
	return &Add{songlistWidget: songlistWidget, mpdClient: mpdClient}
}

func (cmd *Add) Reset() {
	cmd.song = song.New()
}

func (cmd *Add) Execute(t lexer.Token) error {
	var err error

	switch t.Class {
	case lexer.TokenIdentifier:
		if len(cmd.song.Tags["file"]) > 0 {
			return fmt.Errorf("Cannot add multiple paths on the same command line.")
		}
		cmd.song.Tags["file"] = []rune(t.String())

	case lexer.TokenEnd:
		cursor := false
		songlistWidget := cmd.songlistWidget()

		if len(cmd.song.Tags["file"]) == 0 {
			cmd.song = songlistWidget.CursorSong()
			if cmd.song == nil {
				return fmt.Errorf("No song under cursor, cannot add without any parameters.")
			}
			cursor = true
		}

		client := cmd.mpdClient()
		if client == nil {
			return fmt.Errorf("Cannot play: not connected to MPD")
		}

		err = client.Add(cmd.song.StringTags["file"])

		if cursor && err == nil {
			songlistWidget.MoveCursor(1)
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return err
}
