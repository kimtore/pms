// Search index songs

package song

import (
	"github.com/ambientsound/pms/song"
)

// Song is a Bleve document representing a song.Song object.
type Song struct {
	Album       string
	AlbumArtist string
	Artist      string
	File        string
	Genre       string
	Title       string
}

// New generates a indexable Song document, containing some fields from the song.Song type.
func New(s *song.Song) (is Song) {
	is.Album = s.TagString("album")
	is.AlbumArtist = s.TagString("albumartist")
	is.Artist = s.TagString("artist")
	is.File = s.TagString("file")
	is.Genre = s.TagString("genre")
	is.Title = s.TagString("title")
	return
}
