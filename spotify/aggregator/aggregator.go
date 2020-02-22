package spotify_aggregator

import (
	"fmt"

	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/spotify/library"
	"github.com/ambientsound/pms/spotify/playlists"
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

func ListWithID(client spotify.Client, id string, limit int) (*spotify_tracklist.List, error) {
	sid := spotify.ID(id)

	playlist, err := client.GetPlaylist(sid)
	if err != nil {
		return nil, err
	}

	tracks, err := client.GetPlaylistTracksOpt(sid, &spotify.Options{
		Limit: &limit,
	}, "")
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromPlaylistTrackPage(client, tracks)
	if err != nil {
		return nil, err
	}

	lst.SetName(fmt.Sprintf("%s by %s", playlist.Name, playlist.Owner.DisplayName))
	lst.SetID(id)
	lst.SetVisibleColumns(options.GetList(options.Columns))

	return lst, nil
}

func MyPrivatePlaylists(client spotify.Client, limit int) (*spotify_playlists.List, error) {
	playlists, err := client.CurrentUsersPlaylistsOpt(&spotify.Options{
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}

	lst, err := spotify_playlists.New(client, playlists)
	if err != nil {
		return nil, err
	}

	lst.SetName("My playlists")
	lst.SetID(spotify_library.MyPlaylists)
	lst.SetVisibleColumns(lst.ColumnNames())

	return lst, nil
}

func MyTracks(client spotify.Client, limit int) (*spotify_tracklist.List, error) {
	tracks, err := client.CurrentUsersTracksOpt(&spotify.Options{
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromSavedTrackPage(client, tracks)
	if err != nil {
		return nil, err
	}

	lst.SetName("Saved tracks")
	lst.SetID(spotify_library.MyTracks)
	lst.SetVisibleColumns(options.GetList(options.Columns))
	_ = lst.Sort(options.GetList(options.Sort))

	return lst, nil
}

func TopTracks(client spotify.Client, limit int) (*spotify_tracklist.List, error) {
	tracks, err := client.CurrentUsersTopTracksOpt(&spotify.Options{
		Limit: &limit,
	})
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromFullTrackPage(client, tracks)
	if err != nil {
		return nil, err
	}

	lst.SetName("Top tracks")
	lst.SetID(spotify_library.TopTracks)
	lst.SetVisibleColumns(options.GetList(options.Columns))

	return lst, nil
}

func NewReleases(client spotify.Client) (*spotify_tracklist.List, error) {
	albums, err := client.NewReleases()
	if err != nil {
		return nil, err
	}

	lst, err := spotify_tracklist.NewFromSimpleAlbumPage(client, albums)
	if err != nil {
		return nil, err
	}

	lst.SetName("New releases")
	lst.SetID(spotify_library.NewReleases)
	lst.SetVisibleColumns(options.GetList(options.Columns))

	return lst, nil
}
