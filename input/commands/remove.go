package commands

import (
	"fmt"

	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/widgets"
)

// Remove removes songs from songlists.
type Remove struct {
	songlistWidget func() *widgets.SonglistWidget
}

func NewRemove(songlistWidget func() *widgets.SonglistWidget) *Remove {
	return &Remove{songlistWidget: songlistWidget}
}

func (cmd *Remove) Reset() {
}

func (cmd *Remove) Execute(t lexer.Token) error {
	var err error

	switch t.Class {
	case lexer.TokenEnd:
		songlistWidget := cmd.songlistWidget()
		list := songlistWidget.Songlist()

		if songlistWidget.CursorSong() == nil {
			return fmt.Errorf("No song selected, cannot remove without any parameters.")
		}

		index := songlistWidget.Cursor()
		err = list.Remove(index)

	default:
		return fmt.Errorf("Unknown input '%s', expected END", string(t.Runes))
	}

	return err
}
