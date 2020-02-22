package spotify_aggregator

import (
	spotify_library "github.com/ambientsound/pms/spotify/library"
	spotify_playlists "github.com/ambientsound/pms/spotify/playlists"
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/zmb3/spotify"
)

func Search(client spotify.Client, query string, limit int) (*spotify_tracklist.List, error) {
	results, err := client.SearchOpt(query, spotify.SearchTypeTrack, &spotify.Options{
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}

	return spotify_tracklist.NewFromFullTrackPage(client, results.Tracks)
}

func FeaturedPlaylists(client spotify.Client, limit int) (*spotify_playlists.List, error) {
	message, playlists, err := client.FeaturedPlaylistsOpt(&spotify.PlaylistOptions{
		Options: spotify.Options{
			Limit: &limit,
		},
	})
	if err != nil {
		return nil, err
	}

	lst, err := spotify_playlists.New(client, playlists)
	if err != nil {
		return nil, err
	}

	lst.SetName(message)
	lst.SetID(spotify_library.FeaturedPlaylists)
	lst.SetVisibleColumns(lst.ColumnNames())

	return lst, nil
}
