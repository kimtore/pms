package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/spotify/aggregator"
	"github.com/ambientsound/pms/spotify/devices"
	"github.com/ambientsound/pms/spotify/library"
	"github.com/zmb3/spotify"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

// List navigates and manipulates songlists.
type List struct {
	command
	api       api.API
	client    *spotify.Client
	absolute  int
	duplicate bool
	goto_     bool
	open      bool
	relative  int
	remove    bool
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
		case "open":
			cmd.open = true
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

	case cmd.open:
		row := cmd.api.List().CursorRow()
		if row == nil {
			return fmt.Errorf("no playlist selected")
		}
		return cmd.Goto(row[list.RowIDKey])

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
	case spotify_library.MyPlaylists:
		lst, err = spotify_aggregator.MyPrivatePlaylists(*cmd.client, limit)
	case spotify_library.FeaturedPlaylists:
		lst, err = spotify_aggregator.FeaturedPlaylists(*cmd.client, limit)
	case spotify_library.MyTracks:
		lst, err = spotify_aggregator.MyTracks(*cmd.client, limit)
	case spotify_library.TopTracks:
		lst, err = spotify_aggregator.TopTracks(*cmd.client, limit)
	case spotify_library.Devices:
		lst, err = spotify_devices.New(*cmd.client)
	default:
		lst, err = spotify_aggregator.ListWithID(*cmd.client, id, limit)
		if err != nil {
			break
		}
	}
	dur := time.Since(t)

	if err != nil {
		return err
	}

	log.Debugf("Retrieved %s with %d items in %s", id, lst.Len(), dur.String())
	log.Infof("Loaded %s.", lst.Name())

	// Reset cursor
	lst.SetCursor(0)

	cmd.api.SetList(lst)

	return nil
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
