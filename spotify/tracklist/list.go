package spotify_tracklist

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/list"
	spotify_albums "github.com/ambientsound/pms/spotify/albums"
	"github.com/ambientsound/pms/utils"
	"github.com/zmb3/spotify"
)

type List struct {
	list.Base
	tracks map[string]spotify.FullTrack
}

var _ list.List = &List{}

func NewFromFullTrackPage(client spotify.Client, source *spotify.FullTrackPage) (*List, error) {
	var err error

	tracks := make([]spotify.FullTrack, 0, source.Total)

	for err == nil {
		tracks = append(tracks, source.Tracks...)
		err = client.NextPage(source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil
}

func NewFromSimpleTrackPageAndAlbum(client spotify.Client, source *spotify.SimpleTrackPage, album spotify.SimpleAlbum) (*List, error) {
	var err error

	tracks := make([]spotify.FullTrack, 0, source.Total)

	for err == nil {
		for _, track := range source.Tracks {
			tracks = append(tracks, AlbumTrack(track, album))
		}
		err = client.NextPage(source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil
}

func NewFromSavedTrackPage(client spotify.Client, source *spotify.SavedTrackPage) (*List, error) {
	var err error

	tracks := make([]spotify.FullTrack, 0, source.Total)

	for err == nil {
		for _, track := range source.Tracks {
			tracks = append(tracks, track.FullTrack)
		}
		err = client.NextPage(source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil
}

func NewFromPlaylistTrackPage(client spotify.Client, source *spotify.PlaylistTrackPage) (*List, error) {
	var err error

	tracks := make([]spotify.FullTrack, 0, source.Total)

	for err == nil {
		for _, track := range source.Tracks {
			tracks = append(tracks, track.Track)
		}
		err = client.NextPage(source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil

}

func AlbumTrack(track spotify.SimpleTrack, album spotify.SimpleAlbum) spotify.FullTrack {
	return spotify.FullTrack{
		SimpleTrack: track,
		Album:       album,
	}
}

func NewFromSimpleAlbumPage(client spotify.Client, source *spotify.SimpleAlbumPage) (*List, error) {
	var err error

	tracks := make([]spotify.FullTrack, 0, source.Total)
	albums, err := spotify_albums.NewFromSimpleAlbumPage(client, source)
	if err != nil {
		return nil, err
	}

	for i := 0; i < albums.Len(); i++ {
		album := albums.Album(i)
		t, err := client.GetAlbumTracks(album.ID)
		if err != nil {
			break
		}
		lst, err := NewFromSimpleTrackPageAndAlbum(client, t, *album)
		if err != nil {
			break
		}
		tracks = append(tracks, lst.Tracks()...)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	return NewFromTracks(tracks), nil
}

func NewFromTracks(tracks []spotify.FullTrack) *List {
	this := &List{
		tracks: make(map[string]spotify.FullTrack, len(tracks)),
	}
	this.Clear()
	for _, track := range tracks {
		this.tracks[track.ID.String()] = track
		this.Add(FullTrackRow(track))
	}
	return this
}

func FullTrackRow(track spotify.FullTrack) list.Row {
	return list.Row{
		list.RowIDKey: track.ID.String(),
		"album":       track.Album.Name,
		"albumArtist": strings.Join(artistNames(track.Album.Artists), ", "),
		"artist":      strings.Join(artistNames(track.Artists), ", "),
		"date":        track.Album.ReleaseDateTime().Format("2006-01-02"),
		"time":        utils.TimeString(track.Duration / 1000),
		"title":       track.Name,
		"track":       fmt.Sprintf("%02d", track.TrackNumber),
		"disc":        fmt.Sprintf("%d", track.DiscNumber),
		"popularity":  fmt.Sprintf("%1.2f", float64(track.Popularity)/100),
		"year":        track.Album.ReleaseDateTime().Format("2006"),
	}
}

// CursorSong returns the song currently selected by the cursor.
func (l *List) CursorSong() *spotify.FullTrack {
	return l.Song(l.Cursor())
}

// Song returns the song at a specific index.
func (l *List) Song(index int) *spotify.FullTrack {
	row := l.Row(index)
	if row == nil {
		return nil
	}
	track := l.tracks[row.ID()]
	return &track
}

// Selection returns all the selected songs as a new track list.
func (l *List) Selection() List {
	indices := l.SelectionIndices()
	tracks := make([]spotify.FullTrack, len(indices))

	for i, index := range indices {
		tracks[i] = *l.Song(index)
	}

	return *NewFromTracks(tracks)
}

func (l *List) Tracks() []spotify.FullTrack {
	tracks := make([]spotify.FullTrack, len(l.tracks))
	for i := 0; i < l.Len(); i++ {
		tracks[i] = *l.Song(i)
	}
	return tracks
}
