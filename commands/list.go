package commands

import (
	"fmt"
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/spotify/library"
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/zmb3/spotify"
	"strconv"
	"strings"
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// List navigates and manipulates songlists.
type List struct {
	command
	api       api.API
	client    *spotify.Client
	relative  int
	absolute  int
	duplicate bool
	remove    bool
	goto_     bool
	name      string
}

func NewList(api api.API) Command {
	return &List{
		api:      api,
		absolute: -1,
	}
}

func (cmd *List) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabCompleteVerbs(lit)

	switch tok {
	case lexer.TokenIdentifier:
		switch lit {
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
			cmd.absolute = cmd.api.Db().Len() - 1
		case "goto":
			cmd.goto_ = true
		default:
			i, err := strconv.Atoi(lit)
			if err != nil {
				return fmt.Errorf("cannot navigate lists: position '%s' is not recognized, and is not a number", lit)
			}
			cmd.absolute = i - 1
		}
	default:
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	if cmd.goto_ {
		for tok != lexer.TokenEnd {
			tok, lit = cmd.ScanIgnoreWhitespace()
			cmd.name += lit
		}

		cmd.Unscan()
		cmd.setTabComplete(cmd.name, cmd.api.Db().Keys())
	} else {
		cmd.setTabCompleteEmpty()
	}

	return cmd.ParseEnd()
}

func (cmd *List) Exec() error {
	switch {
	case cmd.goto_:
		return cmd.Goto(cmd.name)

	case cmd.relative != 0:
		cmd.api.Db().MoveCursor(cmd.relative)
		cmd.api.SetList(cmd.api.Db().Current())

	case cmd.absolute >= 0:
		cmd.api.Db().SetCursor(cmd.absolute)
		cmd.api.SetList(cmd.api.Db().Current())

	case cmd.duplicate:
		tracklist := cmd.api.Tracklist()
		if tracklist == nil {
			return fmt.Errorf("only track lists can be duplicated")
		}
		return fmt.Errorf("duplicate is not implemented")

	case cmd.remove:
		return fmt.Errorf("remove is not implemented")
	}

	return nil
}

// Goto loads an external list and applies default columns and sorting.
// Local, cached versions are tried first.
func (cmd *List) Goto(id string) error {
	var err error
	var lst list.List

	// Set Spotify object request limit. Ignore user-defined max limit here,
	// because big queries will always be faster and consume less bandwidth,
	// when requesting all the data.
	const limit = 50

	// Try a cached version of a named list
	lst = cmd.api.Db().List(cmd.name)
	if lst != nil {
		cmd.api.SetList(lst)
		return nil
	}

	// Other named lists need Spotify access
	cmd.client, err = cmd.api.Spotify()
	if err != nil {
		return err
	}

	t := time.Now()
	switch id {
	case spotify_library.MyTracks:
		lst, err = cmd.gotoMyTracks(limit)
	case spotify_library.TopTracks:
		lst, err = cmd.gotoTopTracks(limit)
	default:
		err = fmt.Errorf("no such stored list: %s", id)
	}
	dur := time.Since(t)

	if err != nil {
		return err
	}

	log.Debugf("Retrieved %s with %d tracks in %s", id, lst.Len(), dur.String())
	log.Infof("Loaded %s.", lst.Name())

	// Show default columns for all named lists
	cols := strings.Split(cmd.api.Options().GetString("columns"), ",")
	lst.SetVisibleColumns(cols)

	// Reset cursor
	lst.SetCursor(0)

	cmd.api.SetList(lst)

	return nil
}

func (cmd *List) gotoMyTracks(limit int) (list.List, error) {
	tracks, err := cmd.client.CurrentUsersTracksOpt(&spotify.Options{
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromSavedTrackPage(*cmd.client, tracks)
	if err != nil {
		return nil, err
	}

	lst.SetName("Saved tracks")
	lst.SetID(spotify_library.MyTracks)

	// Apply default sorting
	sort := strings.Split(cmd.api.Options().GetString("sort"), ",")
	_ = lst.Sort(sort)

	return lst, nil
}

func (cmd *List) gotoTopTracks(limit int) (list.List, error) {
	tracks, err := cmd.client.CurrentUsersTopTracksOpt(&spotify.Options{
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromFullTrackPage(*cmd.client, tracks)
	if err != nil {
		return nil, err
	}

	lst.SetName("Top tracks")
	lst.SetID(spotify_library.TopTracks)

	return lst, nil
}

// setTabCompleteVerbs sets the tab complete list to the list of available sub-commands.
func (cmd *List) setTabCompleteVerbs(lit string) {
	cmd.setTabComplete(lit, []string{
		"down",
		"duplicate",
		"end",
		"goto",
		"home",
		"next",
		"prev",
		"previous",
		"remove",
		"up",
	})
}
