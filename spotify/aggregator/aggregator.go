package spotify_aggregator

import (
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/zmb3/spotify"
)

func Search(client spotify.Client, query string) (*spotify_tracklist.List, error) {
	limit := 50

	results, err := client.SearchOpt(query, spotify.SearchTypeTrack, &spotify.Options{
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}

	return spotify_tracklist.New(client, results.Tracks)
}
