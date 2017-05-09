package commands

import (
	"fmt"
	"strconv"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/songlist"
)

// List navigates and manipulates songlists.
type List struct {
	api       API
	relative  int
	absolute  int
	duplicate bool
	remove    bool
}

func NewList(api API) Command {
	return &List{
		api:      api,
		absolute: -1,
	}
}

func (cmd *List) Execute(t lexer.Token) error {
	var err error
	var index int

	s := t.String()
	ui := cmd.api.UI()
	songlistWidget := cmd.api.SonglistWidget()

	switch t.Class {

	case lexer.TokenIdentifier:
		switch s {
		case "duplicate":
			cmd.duplicate = true
		case "remove":
			cmd.remove = true
		case "up", "prev", "previous":
			cmd.relative = -1
		case "down", "next":
			cmd.relative = 1
		case "home":
			cmd.absolute = 0
		case "end":
			cmd.absolute = songlistWidget.SonglistsLen() - 1
		default:
			i, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("Cannot navigate lists: position '%s' is not recognized, and is not a number", s)
			}
			switch {
			case cmd.relative != 0 || cmd.absolute != -1:
				return fmt.Errorf("Only one number allowed when setting list position")
			case cmd.relative != 0:
				cmd.relative *= i
			default:
				cmd.absolute = i - 1
			}
		}

	case lexer.TokenEnd:
		switch {
		case cmd.duplicate:
			console.Log("Duplicating current songlist.")
			orig := songlistWidget.Songlist()
			list := songlist.New()
			err = orig.Duplicate(list)
			if err != nil {
				return fmt.Errorf("Error during songlist duplication: %s", err)
			}
			name := fmt.Sprintf("%s (copy)", orig.Name())
			list.SetName(name)
			songlistWidget.AddSonglist(list)
			index = songlistWidget.SonglistsLen() - 1

		case cmd.remove:
			list := songlistWidget.Songlist()
			console.Log("Removing current songlist '%s'.", list.Name())

			err = list.Delete()
			if err != nil {
				return fmt.Errorf("Cannot remove songlist: %s", err)
			}

			index, err = songlistWidget.SonglistIndex()

			// If we got an error here, it means that the current songlist is
			// not in the list of songlists. In this case, we can reset to the
			// fallback songlist.
			if err != nil {
				fallback := songlistWidget.FallbackSonglist()
				if fallback == nil {
					return fmt.Errorf("No songlists left.")
				}
				console.Log("Songlist was not found in the list of songlists. Activating fallback songlist '%s'.", fallback.Name())
				ui.PostFunc(func() {
					songlistWidget.SetSonglist(fallback)
				})
				return nil
			} else {
				songlistWidget.RemoveSonglist(index)
			}

			// If removing the last songlist, we need to decrease the songlist index by one.
			if index == songlistWidget.SonglistsLen() {
				index--
			}

			console.Log("Removed songlist, now activating songlist no. %d", index)

		case cmd.relative != 0:
			index, err = songlistWidget.SonglistIndex()
			if err != nil {
				index = 0
			}
			index += cmd.relative
			if !songlistWidget.ValidSonglistIndex(index) {
				len := songlistWidget.SonglistsLen()
				index = (index + len) % len
			}
			console.Log("Switching songlist index to relative %d, equalling absolute %d", cmd.relative, index)

		case cmd.absolute >= 0:
			console.Log("Switching songlist index to absolute %d", cmd.absolute)
			index = cmd.absolute

		default:
			return fmt.Errorf("Unexpected END, expected position. Try one of: next prev <number>")
		}

		ui.PostFunc(func() {
			err = songlistWidget.SetSonglistIndex(index)
		})

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return err
}
