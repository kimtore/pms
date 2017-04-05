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
	is.Album = s.TagString("album")
	is.AlbumArtist = s.TagString("albumartist")
	is.Artist = s.TagString("artist")
	is.File = s.TagString("file")
	is.Genre = s.TagString("genre")
	is.Title = s.TagString("title")
	return
}
