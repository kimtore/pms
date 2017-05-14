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

func (cmd *Cursor) parseNextOf() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	if tok != lexer.TokenIdentifier {
		return fmt.Errorf("Unexpected %v, expected END", lit)
	}

	cmd.setNextOfTags(lit)
	//if err == nil {
	//cmd.absolute = cmd.runNextOf(direction)
	//}

	return nil
}

// ParseEnd parses to the end, and returns an error if the end hasn't been reached.
func (cmd *Cursor) ParseEnd() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	if tok != lexer.TokenEnd {
		return fmt.Errorf("Unexpected %v, expected END", lit)
	}
	return nil
}

// Parse parses cursor movement.
func (cmd *Cursor) Parse(s *lexer.Scanner) error {
	// Boilerplate
	cmd.S = s
	songlistWidget := cmd.api.SonglistWidget()
	list := cmd.api.Songlist()

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

func (cmd *Cursor) Execute(class int, s string) error {
	var err error

	songlistWidget := cmd.api.SonglistWidget()
	list := cmd.api.Songlist()

	if cmd.finished && class != lexer.TokenEnd {
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	switch class {

	case lexer.TokenIdentifier:
		switch s {
		case "up":
			cmd.relative = -1
		case "down":
			cmd.relative = 1
		case "pgup":
			fallthrough
		case "pageup":
			_, y := songlistWidget.Size()
			cmd.relative = -y
		case "pgdn":
			fallthrough
		case "pagedn":
			fallthrough
		case "pagedown":
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
			return nil
		case "prev-of":
			cmd.nextOfDirection = -1
			return nil
		default:
			switch cmd.nextOfDirection {
			case 1, -1:
				err = cmd.setNextOfTags(s)
				if err == nil {
					cmd.absolute = cmd.runNextOf(cmd.nextOfDirection)
				}
			default:
				i, err := strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("Cannot move cursor: input '%s' is not recognized, and is not a number", s)
				}
				cmd.relative = i
			}
		}
		cmd.finished = true

	case lexer.TokenEnd:
		switch {
		case !cmd.finished:
			return fmt.Errorf("Unexpected END, expected cursor offset.")

		case cmd.current:
			currentSong := cmd.api.Song()
			if currentSong == nil {
				return fmt.Errorf("No song is currently playing.")
			}
			err = list.CursorToSong(currentSong)

		case cmd.relative != 0:
			list.MoveCursor(cmd.relative)

		default:
			list.SetCursor(cmd.absolute)
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
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
