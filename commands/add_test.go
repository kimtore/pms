package commands_test

import (
	"strings"
	"testing"

	"github.com/ambientsound/pms/commands"
	"github.com/stretchr/testify/assert"
)

var addTests = []commands.Test{
	// Valid forms
	{``, true, initSongTags, nil, []string{}},
	{`foo bar baz`, true, nil, nil, []string{}},
	{`http://example.com/stream.mp3?foo=bar&baz=foo foo bar baz`, true, nil, nil, []string{}},
	{`|`, true, nil, nil, []string{}},
	{`|{}$`, true, nil, nil, []string{}},

	// No invalid forms, all input is accepted
}

func TestAdd(t *testing.T) {
	commands.TestVerb(t, "add", addTests)
}

// FIXME: add this callback to test #3. Not working because Queue doesn't add directly.
func testMultipleSongsAdded(data *commands.TestData) {
	files := strings.Split(data.Test.Input, " ")

	assert.Equal(data.T, len(files), data.Api.Songlist().Len(), "Number of URIs added differs from songlist length")
	for i, song := range data.Api.Songlist().Songs() {
		assert.Equal(data.T, files[i], song.StringTags["file"], "Song %d should have URI '%s'", i, files[i])
	}
}
