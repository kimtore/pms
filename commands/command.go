// Package commands contains all functionality that is triggered by the user,
// either through keyboard bindings or the command-line interface. New commands
// such as 'sort', 'add', etc. must be implemented here.
package commands

import (
	"fmt"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/parser"
	"github.com/ambientsound/pms/spotify/devices"
	"github.com/ambientsound/pms/spotify/library"
	"github.com/ambientsound/pms/spotify/playlists"
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/ambientsound/pms/utils"
	"sort"
)

const (
	DevicesContext   = "devices"
	GlobalContext    = "global"
	LibraryContext   = "library"
	PlaylistsContext = "playlists"
	TracklistContext = "tracklist"
	WindowsContext   = "windows"
)

// contexts are used to bind keyboard commands to a specific area of the program.
// For instance, the <ENTER> key can be bound to `play` in track lists and `print _id` in other lists.
var contexts = []string{
	DevicesContext,
	GlobalContext,
	LibraryContext,
	PlaylistsContext,
	TracklistContext,
	WindowsContext,
}

// Verbs contain mappings from strings to Command constructors.
// Make sure to add commands here when implementing them.
var Verbs = map[string]func(api.API) Command{
	"auth":      NewAuth,
	"bind":      NewBind,
	"columns":   NewColumns,
	"copy":      NewYank,
	"cursor":    NewCursor,
	"cut":       NewCut,
	"device":    NewDevice,
	"inputmode": NewInputMode,
	"isolate":   NewIsolate,
	"list":      NewList,
	"next":      NewNext,
	"paste":     NewPaste,
	"pause":     NewPause,
	"play":      NewPlay,
	"previous":  NewPrevious,
	"prev":      NewPrevious,
	"print":     NewPrint,
	"q":         NewQuit,
	"quit":      NewQuit,
	"redraw":    NewRedraw,
	"seek":      NewSeek,
	"select":    NewSelect,
	"se":        NewSet,
	"set":       NewSet,
	"show":      NewShow,
	"single":    NewSingle,
	"sort":      NewSort,
	"stop":      NewStop,
	"style":     NewStyle,
	"unbind":    NewUnbind,
	"update":    NewUpdate,
	"viewport":  NewViewport,
	"volume":    NewVolume,
	"yank":      NewYank,
	// "list":      NewList,
}

// Command must be implemented by all commands.
type Command interface {
	// Exec executes the AST generated by the command.
	Exec() error

	// SetScanner assigns a scanner to the command.
	// FIXME: move to constructor?
	SetScanner(*lexer.Scanner)

	// Parse and make an abstract syntax tree. This function MUST NOT have any side effects.
	Parse() error

	// TabComplete returns a set of tokens that could possibly be used as the next
	// command parameter.
	TabComplete() []string

	// Scanned returns a slice of tokens that have been scanned using Parse().
	Scanned() []parser.Token
}

// command is the base class for all commands, implementing the parser and tab completion.
type command struct {
	parser.Parser
	cmdline     string
	tabComplete []string
}

// Return an ordered list of which program contexts active right now.
// Local contexts take precedence over global contexts.
func Contexts(a api.API) []string {
	ctx := make([]string, 0, len(contexts))
	lst := a.List()
	switch lst.(type) {
	case *db.List:
		ctx = append(ctx, WindowsContext)
	case *spotify_library.List:
		ctx = append(ctx, LibraryContext)
	case *spotify_tracklist.List:
		ctx = append(ctx, TracklistContext)
	case *spotify_playlists.List:
		ctx = append(ctx, PlaylistsContext)
	case *spotify_devices.List:
		ctx = append(ctx, DevicesContext)
	}
	ctx = append(ctx, GlobalContext)
	return ctx
}

// New returns the Command associated with the given verb.
func New(verb string, a api.API) Command {
	ctor := Verbs[verb]
	if ctor == nil {
		return nil
	}
	return ctor(a)
}

// Keys returns a string slice with all verbs that can be invoked to run a command.
func Keys() []string {
	keys := make(sort.StringSlice, 0, len(Verbs))
	for verb := range Verbs {
		keys = append(keys, verb)
	}
	keys.Sort()
	return keys
}

// setTabComplete defines a string slice that will be used for tab completion
// at the current point in parsing.
func (c *command) setTabComplete(filter string, s []string) {
	c.tabComplete = utils.TokenFilter(filter, s)
}

// setTabCompleteEmpty removes all tab completions.
func (c *command) setTabCompleteEmpty() {
	c.setTabComplete("", []string{})
}

// ParseTags parses a set of tags until the end of the line, and maintains the
// tab complete list according to a specified song.
func (c *command) ParseTags(possibleTags []string) ([]string, error) {
	c.setTabCompleteEmpty()
	tags := make([]string, 0)
	tag := ""

	for {
		tok, lit := c.Scan()

		switch tok {
		case lexer.TokenWhitespace:
			if len(tag) > 0 {
				tags = append(tags, tag)
			}
			tag = ""
		case lexer.TokenEnd, lexer.TokenComment:
			if len(tag) > 0 {
				tags = append(tags, tag)
			}
			if len(tags) == 0 {
				return nil, fmt.Errorf("Unexpected END, expected tag")
			}
			return tags, nil
		default:
			tag += lit
		}

		c.setTabComplete(tag, possibleTags)
	}
}

// ParseContext parses a single identifier and verifies that it is a program context.
func (c *command) ParseContext() (string, error) {
	tok, lit := c.ScanIgnoreWhitespace()
	c.setTabComplete(lit, contexts)

	if tok != lexer.TokenIdentifier {
		return "", fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	for _, ctx := range contexts {
		if lit == ctx {
			c.setTabCompleteEmpty()
			return lit, nil
		}
	}

	return "", fmt.Errorf("unexpected '%s', expected one of %v", lit, contexts)
}

// TabComplete implements Command.TabComplete.
func (c *command) TabComplete() []string {
	if c.tabComplete == nil {
		// FIXME
		return make([]string, 0)
	}
	return c.tabComplete
}
