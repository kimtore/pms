package spotify_albums

import (
	"strings"

	"github.com/ambientsound/pms/list"
	"github.com/zmb3/spotify"
)

type List struct {
	list.Base
	albums map[string]spotify.SimpleAlbum
}

var _ list.List = &List{}

func NewFromSimpleAlbumPage(client spotify.Client, source *spotify.SimpleAlbumPage) (*List, error) {
	var err error

	albums := make([]spotify.SimpleAlbum, 0, source.Total)

	for err == nil {
		albums = append(albums, source.Albums...)
		err = client.NextPage(source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromAlbums(albums), nil
}

func NewFromSavedAlbumPage(client spotify.Client, source *spotify.SavedAlbumPage) (*List, error) {
	var err error

	albums := make([]spotify.SimpleAlbum, 0, source.Total)

	for err == nil {
		for _, album := range source.Albums {
			albums = append(albums, album.SimpleAlbum)
		}
		err = client.NextPage(source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromAlbums(albums), nil
}

func NewFromAlbums(albums []spotify.SimpleAlbum) *List {
	this := &List{
		albums: make(map[string]spotify.SimpleAlbum, len(albums)),
	}
	this.Clear()
	for _, album := range albums {
		this.albums[album.ID.String()] = album
		this.Add(SimpleAlbumRow(album))
	}
	return this
}

func SimpleAlbumRow(album spotify.SimpleAlbum) list.Row {
	return list.Row{
		list.RowIDKey: album.ID.String(),
		"album":       album.Name,
		"albumArtist": strings.Join(artistNames(album.Artists), ", "),
		"artist":      strings.Join(artistNames(album.Artists), ", "),
		"date":        album.ReleaseDateTime().Format("2006-01-02"),
		"title":       album.Name,
		"group":       album.AlbumGroup,
		"type":        album.AlbumType,
		"year":        album.ReleaseDateTime().Format("2006"),
	}
}

// Song returns the song at a specific index.
func (l *List) Album(index int) *spotify.SimpleAlbum {
	row := l.Row(index)
	if row == nil {
		return nil
	}
	album := l.albums[row.ID()]
	return &album
}
