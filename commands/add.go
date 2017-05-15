package commands

import (
	"fmt"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
)

// Add adds songs to MPD's queue.
type Add struct {
	newcommand
	api      api.API
	songlist songlist.Songlist
}

func NewAdd(api api.API) Command {
	return &Add{
		api:      api,
		songlist: songlist.New(),
	}
}

// Parse implements Command.
func (cmd *Add) Parse() error {

	// Add all songs specified
Loop:
	for {
		tok, lit := cmd.ScanIgnoreWhitespace()
		switch tok {
		case lexer.TokenIdentifier:
		case lexer.TokenWhitespace, lexer.TokenEnd:
			break Loop
		default:
			return fmt.Errorf("Unexpected %v, expected identifier", lit)
		}
		addSong := song.New()
		addSong.SetTags(mpd.Attrs{"file": lit})
		cmd.songlist.Add(addSong)
	}

	// No songs specified on command line. Use songlist selection instead.
	if cmd.songlist.Len() == 0 {
		cmd.songlist = cmd.api.Songlist().Selection()
		if cmd.songlist.Len() == 0 {
			return fmt.Errorf("No selection, cannot add without any parameters.")
		}
	}

	return nil
}

// Exec is the next Execute(), evading the old system
func (cmd *Add) Exec() error {
	list := cmd.api.SonglistWidget().Songlist()
	queue := cmd.api.Queue()

	err := queue.AddList(cmd.songlist)
	if err != nil {
		return err
	}

	list.ClearSelection()
	list.MoveCursor(1)
	len := cmd.songlist.Len()
	if len == 1 {
		song := cmd.songlist.Songs()[0]
		cmd.api.Message("Added to queue: %s", song.StringTags["file"])
	} else {
		cmd.api.Message("Added %d songs to queue.", len)
	}

	return nil
}
