package commands_test

import (
	"testing"

	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/song"
)

var cursorTests = []commands.Test{
	// Valid forms
	{`6`, true, nil, nil, []string{}},
	{`+8`, true, nil, nil, []string{}},
	{`-1`, true, nil, nil, []string{}},
	{`up`, true, nil, nil, []string{}},
	{`down`, true, nil, nil, []string{}},
	// FIXME: depends on SonglistWidget, which is not mocked
	//{`pgup`, true},
	//{`pageup`, true},
	//{`pagedn`, true},
	//{`pagedown`, true},
	//{`high`, true},
	//{`middle`, true},
	//{`low`, true},
	{`home`, true, nil, nil, []string{}},
	{`end`, true, nil, nil, []string{}},
	{`current`, true, nil, nil, []string{}},
	{`random`, true, nil, nil, []string{}},
	{`nextOf tag1 tag2`, true, nil, nil, []string{}},
	{`prevOf tag1 tag2`, true, nil, nil, []string{}},

	// Invalid forms
	{`up 1`, false, nil, nil, []string{}},
	{`down 1`, false, nil, nil, []string{}},
	// FIXME: depends on SonglistWidget, which is not mocked
	//{`pgup 1`, false},
	//{`pageup 1`, false},
	//{`pagedn 1`, false},
	//{`pagedown 1`, false},
	//{`high 1`, false},
	//{`middle 1`, false},
	//{`low 1`, false},
	{`home 1`, false, nil, nil, []string{}},
	{`end 1`, false, nil, nil, []string{}},
	{`current 1`, false, nil, nil, []string{}},
	{`random 1`, false, nil, nil, []string{}},
	{`nextOf`, false, nil, nil, []string{}},
	{`nextOf `, false, initSongTags, nil, []string{"artist", "title"}},
	{`nextOf t`, true, initSongTags, nil, []string{"title"}},
	{`prevOf`, false, nil, nil, []string{}},
	{`prevOf `, false, initSongTags, nil, []string{"artist", "title"}},

	// Tab completion
	{``, false, nil, nil, []string{
		"current",
		"down",
		"end",
		"high",
		"home",
		"low",
		"middle",
		"nextOf",
		"pagedn",
		"pagedown",
		"pageup",
		"pgdn",
		"pgup",
		"prevOf",
		"random",
		"up",
	}},
	{`page`, false, nil, nil, []string{
		"pagedn",
		"pagedown",
		"pageup",
	}},
}

func TestCursor(t *testing.T) {
	commands.TestVerb(t, "cursor", cursorTests)
}

func initSongTags(data *commands.TestData) {
	s := song.New()
	s.SetTags(mpd.Attrs{
		"artist": "foo",
		"title":  "bar",
	})
	data.Api.Songlist().Add(s)
}
