// package spotify_library provides access to pre-defined Spotify content
// such as playlists, library, top artists, categories, and so on.
//
// This package also provides access to internal functions such as devices and clipboard.
package spotify_library

import (
	"github.com/ambientsound/pms/list"
)

type List struct {
	list.Base
}

var _ list.List = &List{}

const (
	listName = "description"
)

var rows = []list.Row{
	{
		list.RowIDKey: "my playlists",
		listName:      "Personal playlists from my Spotify library",
	},
}

func New() *List {
	this := &List{}
	this.Clear()
	this.SetID("spotify_library")
	this.SetName("Libraries and discovery")
	this.SetVisibleColumns([]string{listName})
	for _, row := range rows {
		this.Add(row)
	}
	return this
}
