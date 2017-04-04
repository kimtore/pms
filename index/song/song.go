// Search index songs

package song

import (
	"github.com/ambientsound/pms/song"
)

// The Song struct is a simplified representationAll fields will be indexed.
type Song struct {
	Album       string
	AlbumArtist string
	Artist      string
	File        string
	Genre       string
	Title       string
}

// NewIndexedSongDocument generates a indexable Song document based on the Song type.
func New(s *song.Song) (is Song) {
	is.Album = s.Tags["Album"]
	is.AlbumArtist = s.Tags["AlbumArtist"]
	is.Artist = s.Tags["Artist"]
	is.File = s.Tags["file"]
	is.Genre = s.Tags["Genre"]
	is.Title = s.Tags["Title"]
	return
}
