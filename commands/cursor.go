package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/parser"
	"github.com/ambientsound/pms/song"
)

// Cursor moves the cursor in a songlist widget. It can take human-readable
// parameters such as 'up' and 'down', and it also accepts relative positions
// if a number is given.
type Cursor struct {
	parser.Parser
	command
	api             api.API
	absolute        int
	current         bool
	finished        bool
	nextOfDirection int
	nextOfTags      []string
	relative        int
	tabComplete     []string
}

// NewCursor returns Cursor.
func NewCursor(api api.API) Command {
	return &Cursor{
		api: api,
	}
}

// FIXME: remove when Scanned() is removed from 'command' struct, and
// responsibility handed over to parser.Parser
func (cmd *Cursor) Scanned() []parser.Token {
	return cmd.Parser.Scanned()
}

// Execute implements Command.Execute.
// FIXME: boilerplate until Execute is removed from interface
func (cmd *Cursor) Execute(class int, s string) error {
	if class == lexer.TokenEnd {
		return cmd.Exec()
	}
	cmd.cmdline += " " + s
	return nil
}

// filter returns a subset of tokens that match the specified prefix.
func (cmd *Cursor) filter(match string, tokens []string) []string {
	dest := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		if strings.HasPrefix(tok, match) {
			dest = append(dest, tok)
		}
	}
	return dest
}

func (cmd *Cursor) TabComplete() []string {
	return cmd.tabComplete
}

func (cmd *Cursor) tabCompleteEmpty() {
	cmd.tabComplete = []string{}
}

func (cmd *Cursor) tabCompleteVerbs(lit string) {
	cmd.tabComplete = cmd.filter(lit, []string{
		"current",
		"down",
		"end",
		"home",
		"next-of",
		"pagedn",
		"pagedown",
		"pageup",
		"pgdn",
		"pgup",
		"prev-of",
		"random",
		"up",
	})
}

func (cmd *Cursor) tabCompleteSong(lit string, song *song.Song) {
	if song == nil {
		cmd.tabCompleteEmpty()
		return
	}
	cmd.tabComplete = cmd.filter(lit, song.TagKeys())
}

// Parse parses cursor movement.
func (cmd *Cursor) Parse(s *lexer.Scanner) error {
	songlistWidget := cmd.api.SonglistWidget()
	list := cmd.api.Songlist()

	// FIXME: initial verb scan boilerplate?
	cmd.S = s

	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.tabCompleteVerbs(lit)
	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("Unexpected %v, expected identifier", lit)
	}

	switch lit {
	case "up":
		cmd.relative = -1
	case "down":
		cmd.relative = 1
	case "pgup", "pageup":
		_, y := songlistWidget.Size()
		cmd.relative = -y
	case "pgdn", "pagedn", "pagedown":
		_, y := songlistWidget.Size()
		cmd.relative = y
	case "home":
		cmd.absolute = 0
	case "end":
		cmd.absolute = list.Len() - 1
	case "current":
		cmd.current = true
	case "random":
		cmd.absolute = cmd.random()
	case "next-of":
		cmd.nextOfDirection = 1
		return cmd.parseNextOf()
	case "prev-of":
		cmd.nextOfDirection = -1
		return cmd.parseNextOf()
	default:
		i, err := strconv.Atoi(lit)
		if err != nil {
			return fmt.Errorf("Cursor command '%s' not recognized, and is not a number", lit)
		}
		cmd.relative = i
	}

	cmd.tabCompleteEmpty()

	return cmd.ParseEnd()
}

// Exec is the next Execute(), evading the old system
func (cmd *Cursor) Exec() error {
	// FIXME: move this code out of Command
	reader := strings.NewReader(cmd.cmdline)
	scanner := lexer.NewScanner(reader)
	err := cmd.Parse(scanner)
	if err != nil {
		return err
	}

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

// parseNextOf parses a set of tags.
func (cmd *Cursor) parseNextOf() error {
	song := cmd.api.Songlist().CursorSong()
	cmd.tabCompleteEmpty()

	for {
		tok, lit := cmd.Scan()

		switch tok {
		case lexer.TokenWhitespace:
			cmd.tabCompleteSong("", song)
		case lexer.TokenIdentifier:
			cmd.tabCompleteSong(lit, song)
			cmd.nextOfTags = append(cmd.nextOfTags, strings.ToLower(lit))
		case lexer.TokenEnd:
			if len(cmd.nextOfTags) == 0 {
				return fmt.Errorf("Unexpected END, expected tag", lit)
			}
			return nil
		default:
			return fmt.Errorf("Unexpected %v, expected tag", lit)
		}
	}
}

// runNextOf finds the next song with different tags.
func (cmd *Cursor) runNextOf() int {
	list := cmd.api.Songlist()
	index := list.Cursor()
	return list.NextOf(cmd.nextOfTags, index, cmd.nextOfDirection)
}
