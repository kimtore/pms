package spotify_playlists

import (
	"fmt"
	"github.com/ambientsound/pms/list"
	"github.com/zmb3/spotify"
	"strconv"
)

type List struct {
	list.Base
	playlists map[string]spotify.SimplePlaylist
}

var _ list.List = &List{}

func New(client spotify.Client, source *spotify.SimplePlaylistPage) (*List, error) {
	var err error

	playlists := make([]spotify.SimplePlaylist, 0, source.Total)

	for err == nil {
		playlists = append(playlists, source.Playlists...)
		err = client.NextPage(source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromPlaylists(playlists), nil
}

func NewFromPlaylists(playlists []spotify.SimplePlaylist) *List {
	this := &List{
		playlists: make(map[string]spotify.SimplePlaylist, len(playlists)),
	}
	this.Clear()
	for _, playlist := range playlists {
		this.playlists[playlist.ID.String()] = playlist
		this.Add(Row(playlist))
	}
	return this
}

func Row(playlist spotify.SimplePlaylist) list.Row {
	return list.Row{
		list.RowIDKey:   playlist.ID.String(),
		"name":          playlist.Name,
		"tracks":        fmt.Sprintf("%d", playlist.Tracks.Total),
		"owner":         playlist.Owner.DisplayName,
		"public":        strconv.FormatBool(playlist.IsPublic),
		"collaborative": strconv.FormatBool(playlist.Collaborative),
	}
}

// CursorPlaylist returns the playlist currently selected by the cursor.
func (l *List) CursorPlaylist() *spotify.SimplePlaylist {
	return l.Playlist(l.Cursor())
}

// Song returns the song at a specific index.
func (l *List) Playlist(index int) *spotify.SimplePlaylist {
	row := l.Row(index)
	if row == nil {
		return nil
	}
	playlist := l.playlists[row.ID()]
	return &playlist
}

// Selection returns all the selected songs as a new playlist list.
func (l *List) Selection() List {
	indices := l.SelectionIndices()
	playlists := make([]spotify.SimplePlaylist, len(indices))

	for i, index := range indices {
		playlists[i] = *l.Playlist(index)
	}

	return *NewFromPlaylists(playlists)
}

func (l *List) Playlists() []spotify.SimplePlaylist {
	tracks := make([]spotify.SimplePlaylist, len(l.playlists))
	for i := 0; i < l.Len(); i++ {
		tracks[i] = *l.Playlist(i)
	}
	return tracks
}
