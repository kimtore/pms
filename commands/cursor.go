package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/parser"
)

// Cursor moves the cursor in a songlist widget. It can take human-readable
// parameters such as 'up' and 'down', and it also accepts relative positions
// if a number is given.
type Cursor struct {
	parser.Parser
	command
	api             api.API
	relative        int
	absolute        int
	current         bool
	finished        bool
	nextOfDirection int
	nextOfTags      []string
}

func NewCursor(api api.API) Command {
	return &Cursor{
		api: api,
	}
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

// Parse parses cursor movement.
func (cmd *Cursor) Parse(s *lexer.Scanner) error {
	songlistWidget := cmd.api.SonglistWidget()
	list := cmd.api.Songlist()

	// FIXME: initial verb scan boilerplate?
	cmd.S = s
	tok, lit := cmd.ScanIgnoreWhitespace()
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

	return cmd.ParseEnd()
}

// parseNextOf parses a set of tags.
func (cmd *Cursor) parseNextOf() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("Unexpected %v, expected END", lit)
	}
	cmd.setNextOfTags(lit)
	return nil
}

func (cmd *Cursor) Exec() error {
	// Evade old system
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
		cmd.absolute = cmd.runNextOf(cmd.nextOfDirection)
	case cmd.current:
		currentSong := cmd.api.Song()
		if currentSong == nil {
			return fmt.Errorf("No song is currently playing.")
		}
		err = list.CursorToSong(currentSong)
	}

	switch {
	case cmd.relative != 0:
		list.MoveCursor(cmd.relative)
	default:
		list.SetCursor(cmd.absolute)
	}

	return nil
}

func (cmd *Cursor) random() int {
	len := cmd.api.Songlist().Len()
	if len == 0 {
		return cmd.absolute
	}
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return r.Int() % len
}

func (cmd *Cursor) setNextOfTags(taglist string) error {
	if len(cmd.nextOfTags) != 0 {
		return fmt.Errorf("Unexpected tags, expected END")
	}
	cmd.nextOfTags = strings.Split(strings.ToLower(taglist), ",")
	return nil
}

func (cmd *Cursor) runNextOf(direction int) int {
	list := cmd.api.Songlist()
	len := list.Len()

	offset := func(i int) int {
		if direction > 0 || i == 0 {
			return 0
		}
		return 1
	}

	index := list.Cursor()
	index -= offset(index)
	check := list.Song(index)

	for ; index < len && index >= 0; index += direction {
		song := list.Song(index)
		if song == nil {
			console.Log("NextOf: empty song, break")
			break
		}
		for _, tag := range cmd.nextOfTags {
			if check.StringTags[tag] != song.StringTags[tag] {
				console.Log("NextOf: tag '%s' on source '%s' differs from destination '%s', breaking", tag, check.StringTags[tag], song.StringTags[tag])
				return index + offset(index)
			}
		}
	}

	console.Log("NextOf: fallthrough")
	return index + offset(index)
}
