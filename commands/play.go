package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
)

// Play plays songs in the MPD playlist.
type Play struct {
	api  API
	song *song.Song
	id   int
	pos  int
}

func NewPlay(api API) Command {
	return &Play{
		api: api,
		pos: -1,
	}
}

func (cmd *Play) Execute(t lexer.Token) error {
	var err error

	s := t.String()
	songlistWidget := cmd.api.SonglistWidget()

	switch t.Class {
	case lexer.TokenIdentifier:
		switch s {
		case "cursor":

			cmd.song = songlistWidget.CursorSong()
			if cmd.song == nil {
				return fmt.Errorf("Cannot play: no song under cursor")
			}
		default:
			return nil
		}

	case lexer.TokenEnd:
		client := cmd.api.MpdClient()
		if client == nil {
			return fmt.Errorf("Cannot play: not connected to MPD")
		}

		if cmd.song == nil {
			err = client.Play(-1)
			return err
		}

		id := cmd.song.ID

		// Add song to queue only if we are not operating on the queue
		list := songlistWidget.Songlist()

		if !songlist.IsQueue(list) {
			id, err = client.AddID(cmd.song.StringTags["file"], -1)
			if err != nil {
				return err
			}
		}

		err = client.PlayID(id)
		return err

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return nil
}
