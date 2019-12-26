package spotify_tracklist

import (
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/utils"
	"github.com/zmb3/spotify"
)

type List struct {
	list.Base
	tracks []spotify.FullTrack
}

var _ list.List = &List{}

func New(client spotify.Client, source *spotify.FullTrackPage) (*List, error) {
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

func NewFromTracks(tracks []spotify.FullTrack) *List {
	this := &List{
		tracks: tracks,
	}
	this.Clear()
	for _, track := range this.tracks {
		this.Add(Row(track))
	}
	return this
}

func Row(track spotify.FullTrack) list.Row {
	return list.Row{
		"album":  track.Album.Name,
		"artist": track.Artists[0].Name,
		"time":   utils.TimeString(track.Duration / 1000),
		"title":  track.Name,
	}
}

// CursorSong returns the song currently selected by the cursor.
func (l *List) CursorSong() *spotify.FullTrack {
	return l.Song(l.Cursor())
}

// Song returns the song at a specific index.
func (l *List) Song(index int) *spotify.FullTrack {
	if !l.InRange(index) {
		return nil
	}
	return &l.tracks[index]
}

// Selection returns all the selected songs as a new track list.
func (l *List) Selection() List {
	indices := l.SelectionIndices()
	tracks := make([]spotify.FullTrack, len(indices))

	for i, index := range indices {
		tracks[i] = l.tracks[index]
	}

	return *NewFromTracks(tracks)
}

func (l *List) Tracks() []spotify.FullTrack {
	return l.tracks
}
