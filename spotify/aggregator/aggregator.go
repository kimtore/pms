package spotify_aggregator

import (
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/zmb3/spotify"
)

func Search(client spotify.Client, query string) (*spotify_tracklist.List, error) {
	results, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		return nil, err
	}

	return spotify_tracklist.New(client, *results.Tracks)
}
