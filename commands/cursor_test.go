package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var cursorTests = []commands.CommandTest{
	// Valid forms
	{`6`, true, []string{}},
	{`+8`, true, []string{}},
	{`-1`, true, []string{}},
	{`up`, true, []string{}},
	{`down`, true, []string{}},
	// FIXME: depends on SonglistWidget, which is not mocked
	//{`pgup`, true},
	//{`pageup`, true},
	//{`pagedn`, true},
	//{`pagedown`, true},
	{`home`, true, []string{}},
	{`end`, true, []string{}},
	{`current`, true, []string{}},
	{`random`, true, []string{}},
	{`next-of tag1 tag2`, true, []string{}},
	{`prev-of tag1 tag2`, true, []string{}},

	// Invalid forms
	{`up 1`, false, []string{}},
	{`down 1`, false, []string{}},
	// FIXME: depends on SonglistWidget, which is not mocked
	//{`pgup 1`, false},
	//{`pageup 1`, false},
	//{`pagedn 1`, false},
	//{`pagedown 1`, false},
	{`home 1`, false, []string{}},
	{`end 1`, false, []string{}},
	{`current 1`, false, []string{}},
	{`random 1`, false, []string{}},
	{`next-of`, false, []string{}},
	{`next-of `, false, []string{"artist", "title"}},
	{`next-of t`, true, []string{"title"}},
	{`prev-of`, false, []string{}},
	{`prev-of `, false, []string{"artist", "title"}},

	// Tab completion
	{``, false, []string{
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
	{`page`, false, []string{
		"pagedn",
		"pagedown",
		"pageup",
	}},
}

func TestCursor(t *testing.T) {
	commands.TestVerb(t, "cursor", cursorTests)
}
