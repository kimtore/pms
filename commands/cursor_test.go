package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var cursorTests = []commands.CommandTest{
	// Valid forms
	{`6`, true, nil, []string{}},
	{`+8`, true, nil, []string{}},
	{`-1`, true, nil, []string{}},
	{`up`, true, nil, []string{}},
	{`down`, true, nil, []string{}},
	// FIXME: depends on SonglistWidget, which is not mocked
	//{`pgup`, true},
	//{`pageup`, true},
	//{`pagedn`, true},
	//{`pagedown`, true},
	{`home`, true, nil, []string{}},
	{`end`, true, nil, []string{}},
	{`current`, true, nil, []string{}},
	{`random`, true, nil, []string{}},
	{`next-of tag1 tag2`, true, nil, []string{}},
	{`prev-of tag1 tag2`, true, nil, []string{}},

	// Invalid forms
	{`up 1`, false, nil, []string{}},
	{`down 1`, false, nil, []string{}},
	// FIXME: depends on SonglistWidget, which is not mocked
	//{`pgup 1`, false},
	//{`pageup 1`, false},
	//{`pagedn 1`, false},
	//{`pagedown 1`, false},
	{`home 1`, false, nil, []string{}},
	{`end 1`, false, nil, []string{}},
	{`current 1`, false, nil, []string{}},
	{`random 1`, false, nil, []string{}},
	{`next-of`, false, nil, []string{}},
	{`next-of `, false, nil, []string{"artist", "title"}},
	{`next-of t`, true, nil, []string{"title"}},
	{`prev-of`, false, nil, []string{}},
	{`prev-of `, false, nil, []string{"artist", "title"}},

	// Tab completion
	{``, false, nil, []string{
		"current",
		"down",
		"end",
		"home",
		"next-of",
		"pagedn",
		"pagedown",
		"pageup",
		"pgdn",
		"pgup",
		"prev-of",
		"random",
		"up",
	}},
	{`page`, false, nil, []string{
		"pagedn",
		"pagedown",
		"pageup",
	}},
}

func TestCursor(t *testing.T) {
	commands.TestVerb(t, "cursor", cursorTests)
}
