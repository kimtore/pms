package spotify_list

import (
	"github.com/ambientsound/pms/list"
	"github.com/zmb3/spotify"
)

type List struct {
	list.Base
	columns  []list.Column
	playlist spotify.SimplePlaylist
	tracks   []spotify.PlaylistTrack
}

var _ list.List = &List{}

func (l *List) Tracks() []spotify.PlaylistTrack {
	return l.tracks
}

func Row(track spotify.PlaylistTrack) list.Row {
	return list.Row{
		"title": track.Track.String(),
	}
}

func FromSimplePlaylist(client spotify.Client, source spotify.SimplePlaylist) (*List, error) {
	this := &List{
		columns:  make([]list.Column, 1),
		playlist: source,
		tracks:   make([]spotify.PlaylistTrack, 0),
	}
	tracks, err := client.GetPlaylistTracks(source.ID)

	for err == nil {
		this.tracks = append(this.tracks, tracks.Tracks...)
		err = client.NextPage(tracks)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	for _, track := range this.tracks {
		this.Add(Row(track))
	}

	return this, nil
}
