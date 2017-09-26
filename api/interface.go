package api

import (
	"github.com/ambientsound/pms/songlist"
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

type SonglistWidget interface {
	Bottom() int
	Scroll(int, bool)
	Size() (int, int)
	Top() int
}

type MultibarWidget interface {
	Mode() int
	SetMode(int) error
}

type UI interface {
	PostFunc(func())
	Refresh()
}

type Buffer interface {
	String() string
	Cursor() int
}
