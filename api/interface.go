package api

import (
	"github.com/ambientsound/pms/songlist"
)

type Window int

const (
	WindowLogs Window = iota
	WindowMusic
	WindowPlaylists
)

type Collection interface {
	Activate(songlist.Songlist)
	ActivateIndex(int) error
	Add(songlist.Songlist)
	Current() songlist.Songlist
	Index() (int, error)
	Last() songlist.Songlist
	Len() int
	Remove(int) error
	ValidIndex(int) bool
}

type TableWidget interface {
	GetVisibleBoundaries() (int, int)
	ScrollViewport(int, bool)
	Size() (int, int)
	PositionReadout() string
}

type UI interface {
	ActivateWindow(Window)
	Refresh()
	TableWidget() TableWidget
}
