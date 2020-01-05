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

const (
	Categories          = "categories"
	Devices             = "devices"
	FeaturedPlaylists   = "featured-playlists"
	FollowedArtists     = "followed-artists"
	MyAlbums            = "my-albums"
	MyFollowedPlaylists = "my-followed-playlists"
	MyPlaylists         = "my-playlists"
	MyPrivatePlaylists  = "my-private-playlists"
	MyTracks            = "my-tracks"
	NewReleases         = "new-releases"
	TopArtists          = "top-artists"
	TopTracks           = "top-tracks"
)

var rows = []list.Row{
	{
		list.RowIDKey: Devices,
		listName:      "Player devices",
	},
	{
		list.RowIDKey: MyPrivatePlaylists,
		listName:      "Personal playlists from my Spotify library",
	},
	{
		list.RowIDKey: MyTracks,
		listName:      "All liked songs from my library",
	},
	{
		list.RowIDKey: TopTracks,
		listName:      "Top tracks from my listening history",
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
