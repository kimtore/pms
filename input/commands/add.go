package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/song"

	"github.com/fhs/gompd/mpd"
)

// Add adds songs to MPD's queue.
type Add struct {
	mpdClient func() *mpd.Client
	song      *song.Song
}

func NewAdd(mpdClient func() *mpd.Client) *Add {
	return &Add{mpdClient: mpdClient}
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
		if len(cmd.song.Tags["file"]) == 0 {
			return fmt.Errorf("Unexpected END; expected path to add")
		}

		client := cmd.mpdClient()
		if client == nil {
			return fmt.Errorf("Cannot play: not connected to MPD")
		}

		err = client.Add(cmd.song.TagString("file"))

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return err
}
