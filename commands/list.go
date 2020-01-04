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

func (cmd *List) Goto(id string) error {
	var err error
	var lst list.List

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

	switch id {
	case spotify_library.MyTracks:
		lst, err = cmd.gotoMyTracks()
	default:
		err = fmt.Errorf("no such stored list: %s", id)
	}

	if err != nil {
		return err
	}

	// Show default columns for all named lists
	cols := strings.Split(cmd.api.Options().GetString("columns"), ",")
	lst.SetVisibleColumns(cols)

	// Apply default sorting for all named lists
	sort := strings.Split(cmd.api.Options().GetString("sort"), ",")
	err = lst.Sort(sort)
	if err != nil {
		log.Errorf("unable to sort: %s", err)
	}

	cmd.api.SetList(lst)

	return nil
}

func (cmd *List) gotoMyTracks() (list.List, error) {
	tracks, err := cmd.client.CurrentUsersTracks()
	if err != nil {
		return nil, err
	}

	log.Infof("Loading saved tracks...")

	lst, err := spotify_tracklist.NewFromSavedTrackPage(*cmd.client, tracks)
	if err != nil {
		return nil, err
	}

	log.Debugf("Loaded saved tracks.")

	lst.SetName("Saved tracks")
	lst.SetID(spotify_library.MyTracks)

	return lst, nil
}

func (cmd *List) Exec() error {
	if cmd.goto_ {
		return cmd.Goto(cmd.name)
	}

	return nil
	/*
	switch {
	case cmd.duplicate:
		console.Log("Duplicating current songlist.")
		orig := collection.Current()
		list := songlist.New()
		err = orig.Duplicate(list)
		if err != nil {
			return fmt.Errorf("Error during songlist duplication: %s", err)
		}
		name := fmt.Sprintf("%s (copy)", orig.Name())
		list.SetName(name)
		collection.Add(list)
		index = collection.Len() - 1

	case cmd.remove:
		list := collection.Current()
		console.Log("Removing current songlist '%s'.", list.Name())

		err = list.Delete()
		if err != nil {
			return fmt.Errorf("Cannot remove songlist: %s", err)
		}

		index, err = collection.Index()

		// If we got an error here, it means that the current songlist is
		// not in the list of songlists. In this case, we can reset to the
		// last used songlist.
		if err != nil {
			fallback := collection.Last()
			if fallback == nil {
				return fmt.Errorf("No songlists left.")
			}
			console.Log("Songlist was not found in the list of songlists. Activating fallback songlist '%s'.", fallback.Name())
			// ui.PostFunc(func() { //FIXME???
			collection.Activate(fallback)
			// })
			return nil
		} else {
			collection.Remove(index)
		}

		// If removing the last songlist, we need to decrease the songlist index by one.
		if index == collection.Len() {
			index--
		}

		console.Log("Removed songlist, now activating songlist no. %d", index)

	case cmd.relative != 0:
		index, err = collection.Index()
		if err != nil {
			index = 0
		}
		index += cmd.relative
		if !collection.ValidIndex(index) {
			len := collection.Len()
			index = (index + len) % len
		}
		console.Log("Switching songlist index to relative %d, equalling absolute %d", cmd.relative, index)

	case cmd.absolute >= 0:
		console.Log("Switching songlist index to absolute %d", cmd.absolute)
		index = cmd.absolute

	default:
		return fmt.Errorf("Unexpected END, expected position. Try one of: next prev <number>")
	}

	return nil

	// ui.PostFunc(func() {//FIXME???
	err = collection.ActivateIndex(index)
	*/
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
