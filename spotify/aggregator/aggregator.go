package spotify_aggregator

import (
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
