package commands

import (
	"fmt"
	"strings"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/parser"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
)

// Add adds songs to MPD's queue.
type Add struct {
	parser.Parser
	command
	api      api.API
	songlist songlist.Songlist
}

func NewAdd(api api.API) Command {
	return &Add{
		api:      api,
		songlist: songlist.New(),
	}
}

// FIXME: remove when Scanned() is removed from 'command' struct, and
// responsibility handed over to parser.Parser
func (cmd *Add) Scanned() []parser.Token {
	return cmd.Parser.Scanned()
}

// Execute implements Command.Execute.
// FIXME: boilerplate until Execute is removed from interface
func (cmd *Add) Execute(class int, s string) error {
	if class == lexer.TokenEnd {
		return cmd.Exec()
	}
	cmd.cmdline += " " + s
	return nil
}

// Exec is the next Execute(), evading the old system
func (cmd *Add) Exec() error {
	// FIXME: move this code out of Command
	reader := strings.NewReader(cmd.cmdline)
	scanner := lexer.NewScanner(reader)
	err := cmd.Parse(scanner)
	if err != nil {
		return err
	}

	list := cmd.api.SonglistWidget().Songlist()
	queue := cmd.api.Queue()

	err = queue.AddList(cmd.songlist)
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

// Parse implements Command.
func (cmd *Add) Parse(s *lexer.Scanner) error {

	// FIXME: initial verb scan boilerplate?
	cmd.S = s

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
