// Search index songs

package song

import (
	"github.com/ambientsound/pms/song"
)

// Song is a Bleve document representing a song.Song object.
type Song struct {
	Album       string
	Albumartist string
	Artist      string
	File        string
	Genre       string
	Title       string
	Year        string
}

// New generates a indexable Song document, containing some fields from the song.Song type.
func New(s *song.Song) (is Song) {
	is.Album = s.StringTags["album"]
	is.Albumartist = s.StringTags["albumartist"]
	is.Artist = s.StringTags["artist"]
	is.File = s.StringTags["file"]
	is.Genre = s.StringTags["genre"]
	is.Title = s.StringTags["title"]
	is.Year = s.StringTags["year"]
	return
}
