package api

import (
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
)

type SonglistWidget interface {
	AddSonglist(songlist.Songlist)
	ClearSelection()
	Cursor() int
	CursorSong() *song.Song
	CursorToSong(*song.Song) error
	DisableVisualSelection()
	FallbackSonglist() songlist.Songlist
	Len() int
	List() *list.List
	MoveCursor(int)
	RemoveSonglist(int) error
	Selection() songlist.Songlist
	SetCursor(int)
	SetSonglist(songlist.Songlist)
	SetSonglistIndex(int) error
	Size() (int, int)
	Songlist() songlist.Songlist
	SonglistIndex() (int, error)
	SonglistsLen() int
	ValidSonglistIndex(int) bool
}

type MultibarWidget interface {
	Mode() int
	SetMode(int) error
}

type UI interface {
	PostFunc(func())
	Refresh()
}
