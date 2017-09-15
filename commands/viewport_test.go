package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var viewportTests = []commands.Test{
	// Valid forms
	{`up`, true, nil, nil, []string{}},
	{`down`, true, nil, nil, []string{}},
	// FIXME: depends on SonglistWidget, which is not mocked
	//{`pgup`, true},
	//{`pgdn`, true},
	//{`pageup`, true},
	//{`pagedn`, true},
	//{`pagedown`, true},
	//{`halfpgup`, true},
	//{`halfpgdn`, true},
	//{`halfpageup`, true},
	//{`halfpagedn`, true},
	//{`halfpagedown`, true},
	//{`high`, true},
	//{`middle`, true},
	//{`low`, true},

	// Invalid forms
	{`up 1`, false, nil, nil, []string{}},
	{`down 1`, false, nil, nil, []string{}},
	// FIXME: depends on SonglistWidget, which is not mocked
	//{`pgup 1`, false},
	//{`pgdn 1`, false},
	//{`pageup 1`, false},
	//{`pagedn 1`, false},
	//{`pagedown 1`, false},
	//{`halfpgup 1`, false},
	//{`halfpgdn 1`, false},
	//{`halfpageup 1`, false},
	//{`halfpagedn 1`, false},
	//{`halfpagedown 1`, false},
	//{`high 1`, false},
	//{`middle 1`, false},
	//{`low 1`, false},
	{`nonsense`, false, nil, nil, []string{}},

	// Tab completion
	{``, false, nil, nil, []string{
		"down",
		"halfpagedn",
		"halfpagedown",
		"halfpageup",
		"halfpgdn",
		"halfpgup",
		"high",
		"low",
		"middle",
		"pagedn",
		"pagedown",
		"pageup",
		"pgdn",
		"pgup",
		"up",
	}},
	{`u`, false, nil, nil, []string{
		"up",
	}},
	{`do`, false, nil, nil, []string{
		"down",
	}},
	{`page`, false, nil, nil, []string{
		"pagedn",
		"pagedown",
		"pageup",
	}},
	{`halfpage`, false, nil, nil, []string{
		"halfpagedn",
		"halfpagedown",
		"halfpageup",
	}},
}

func TestViewport(t *testing.T) {
	commands.TestVerb(t, "viewport", viewportTests)
}
