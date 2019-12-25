package spotify_tracklist

import (
	"github.com/ambientsound/pms/list"
	"github.com/zmb3/spotify"
)

type List struct {
	list.Base
	tracks   []spotify.FullTrack
}

var _ list.List = &List{}

func (l *List) Tracks() []spotify.FullTrack {
	return l.tracks
}

func Row(track spotify.FullTrack) list.Row {
	return list.Row{
		"title": track.String(),
	}
}

func New(client spotify.Client, source spotify.FullTrackPage) (*List, error) {
	var err error

	this := &List{
		tracks:   make([]spotify.FullTrack, 0),
	}
	this.Clear()

	for err == nil {
		this.tracks = append(this.tracks, source.Tracks...)
		err = client.NextPage(source)
	}

	if err != spotify.ErrNoMorePages {
		return nil, err
	}

	for _, track := range this.tracks {
		this.Add(Row(track))
	}

	return this, nil
}
