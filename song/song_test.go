package song_test

import (
	"testing"

	"github.com/ambientsound/pms/song"
	"github.com/stretchr/testify/assert"
)

var autofillTests = []struct {
	OriginalSong   song.Song
	AutoFilledSong song.Song
}{
	{
		song.Song{
			Tags: song.Taglist{
				"id":           []rune("1337"),
				"pos":          []rune("33"),
				"time":         []rune("203"),
				"date":         []rune("1986-04-22"),
				"originaldate": []rune("1986-04-22"),
			},
			StringTags: song.StringTaglist{
				"id":           "1337",
				"pos":          "33",
				"time":         "203",
				"date":         "1986-04-22",
				"originaldate": "1986-04-22",
			},
		},
		song.Song{
			ID:       1337,
			Position: 33,
			Time:     203,
			Tags: song.Taglist{
				"id":           []rune("1337"),
				"pos":          []rune("33"),
				"time":         []rune("03:23"),
				"date":         []rune("1986-04-22"),
				"year":         []rune("1986"),
				"originaldate": []rune("1986-04-22"),
				"originalyear": []rune("1986"),
			},
			StringTags: song.StringTaglist{
				"id":           "1337",
				"pos":          "33",
				"time":         "203",
				"date":         "1986-04-22",
				"year":         "1986",
				"originaldate": "1986-04-22",
				"originalyear": "1986",
			},
		},
	},
}

func TestAutofill(t *testing.T) {
	assert := assert.New(t)
	for _, test := range autofillTests {
		test.OriginalSong.AutoFill()
		assert.Equal(test.AutoFilledSong, test.OriginalSong)
	}

}
