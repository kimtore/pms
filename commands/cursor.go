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
)

// Cursor moves the cursor in a songlist widget. It can take human-readable
// parameters such as 'up' and 'down', and it also accepts relative positions
// if a number is given.
type Cursor struct {
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

func (cmd *Cursor) Execute(class int, s string) error {
	var err error

	songlistWidget := cmd.api.SonglistWidget()

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
			cmd.absolute = songlistWidget.Len() - 1
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
					cmd.absolute = cmd.runNextOf()
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
			err = songlistWidget.CursorToSong(currentSong)

		case cmd.relative != 0:
			songlistWidget.MoveCursor(cmd.relative)

		default:
			songlistWidget.SetCursor(cmd.absolute)
		}

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}

func (cmd *Cursor) random() int {
	len := cmd.api.SonglistWidget().Len()
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

func (cmd *Cursor) runNextOf() int {
	songlistWidget := cmd.api.SonglistWidget()
	list := songlistWidget.Songlist()
	len := songlistWidget.Len()

	offset := func(i int) int {
		if cmd.nextOfDirection > 0 || i == 0 {
			return 0
		}
		return 1
	}

	index := songlistWidget.Cursor()
	index -= offset(index)
	check := list.Song(index)

	for ; index < len && index >= 0; index += cmd.nextOfDirection {
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
