package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/widgets"

	"github.com/ambientsound/gompd/mpd"
)

// Add adds songs to MPD's queue.
type Add struct {
	songlistWidget func() *widgets.SonglistWidget
	mpdClient      func() *mpd.Client
	song           *song.Song
	songlist       songlist.Songlist
}

func NewAdd(songlistWidget func() *widgets.SonglistWidget, mpdClient func() *mpd.Client) *Add {
	return &Add{songlistWidget: songlistWidget, mpdClient: mpdClient}
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
		cmd.song.Tags["file"] = []rune(t.String())

	case lexer.TokenEnd:
		songlistWidget := cmd.songlistWidget()

		client := cmd.mpdClient()
		if client == nil {
			return fmt.Errorf("Cannot play: not connected to MPD")
		}

		switch {
		case cmd.song == nil:
			selection := songlistWidget.Selection()
			if selection.Len() == 0 {
				return fmt.Errorf("No selection, cannot add without any parameters.")
			}
			commandList := client.BeginCommandList()
			if commandList == nil {
				return fmt.Errorf("MPD error: cannot begin command list")
			}
			songs := selection.Songs()
			for _, song := range songs {
				commandList.Add(song.StringTags["file"])
			}
			err = commandList.End()
			if err != nil {
				break
			}
			if err == nil {
				songlistWidget.DisableVisualSelection()
				songlistWidget.MoveCursor(1)
			}
		default:
			err = client.Add(cmd.song.StringTags["file"])
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return err
}
