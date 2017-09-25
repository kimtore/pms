package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// Cursor moves the cursor in a songlist widget. It can take human-readable
// parameters such as 'up' and 'down', and it also accepts relative positions
// if a number is given.
type Cursor struct {
	newcommand
	api             api.API
	absolute        int
	current         bool
	finished        bool
	nextOfDirection int
	nextOfTags      []string
	relative        int
}

// NewCursor returns Cursor.
func NewCursor(api api.API) Command {
	return &Cursor{
		api: api,
	}
}

// Parse parses cursor movement.
func (cmd *Cursor) Parse() error {
	songlistWidget := cmd.api.SonglistWidget()
	list := cmd.api.Songlist()

	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteVerbs(lit)

	switch tok {
	// In case of a number, scan the actual number and return
	case lexer.TokenMinus, lexer.TokenPlus:
		cmd.setTabCompleteEmpty()
		cmd.Unscan()
		_, lit, absolute, err := cmd.ParseInt()
		if err != nil {
			return err
		}
		if absolute {
			cmd.absolute = lit
		} else {
			cmd.relative = lit
		}
		return cmd.ParseEnd()

	case lexer.TokenIdentifier:
	default:
		return fmt.Errorf("Unexpected '%v', expected number or identifier", lit)
	}

	switch lit {
	case "up":
		cmd.relative = -1
	case "down":
		cmd.relative = 1
	case "home":
		cmd.absolute = 0
	case "end":
		cmd.absolute = list.Len() - 1
	case "high":
		cmd.absolute = songlistWidget.Top()
	case "middle":
		cmd.absolute = (songlistWidget.Top() + songlistWidget.Bottom()) / 2
	case "low":
		cmd.absolute = songlistWidget.Bottom()
	case "current":
		cmd.current = true
	case "random":
		cmd.absolute = cmd.random()
	case "nextOf":
		cmd.nextOfDirection = 1
		return cmd.parseNextOf()
	case "prevOf":
		cmd.nextOfDirection = -1
		return cmd.parseNextOf()
	default:
		i, err := strconv.Atoi(lit)
		if err != nil {
			return fmt.Errorf("Cursor command '%s' not recognized, and is not a number", lit)
		}
		cmd.relative = i
	}

	cmd.setTabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec is the next Execute(), evading the old system
func (cmd *Cursor) Exec() error {
	list := cmd.api.Songlist()

	switch {
	case cmd.nextOfDirection != 0:
		cmd.absolute = cmd.runNextOf()
	case cmd.current:
		currentSong := cmd.api.Song()
		if currentSong == nil {
			return fmt.Errorf("No song is currently playing.")
		}
		return list.CursorToSong(currentSong)
	}

	switch {
	case cmd.relative != 0:
		list.MoveCursor(cmd.relative)
	default:
		list.SetCursor(cmd.absolute)
	}

	return nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *Cursor) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"current",
		"down",
		"end",
		"high",
		"home",
		"low",
		"middle",
		"nextOf",
		"prevOf",
		"random",
		"up",
	})
}

// random returns a random list index in the songlist.
func (cmd *Cursor) random() int {
	len := cmd.api.Songlist().Len()
	if len == 0 {
		return cmd.absolute
	}
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return r.Int() % len
}

// parseNextOf assigns the nextOf tags and directions, or returns an error if
// no tags are specified.
func (cmd *Cursor) parseNextOf() error {
	var err error
	song := cmd.api.Songlist().CursorSong()
	cmd.nextOfTags, err = cmd.ParseTags(song)
	return err
}

// runNextOf finds the next song with different tags.
func (cmd *Cursor) runNextOf() int {
	list := cmd.api.Songlist()
	index := list.Cursor()
	return list.NextOf(cmd.nextOfTags, index, cmd.nextOfDirection)
}
