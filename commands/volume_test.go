package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var volumeTests = []commands.Test{
	// Valid forms
	{`-2`, true, nil, nil, []string{}},
	{`+13`, true, nil, nil, []string{}},
	{`1329`, true, nil, nil, []string{}},
	{`mute`, true, nil, nil, []string{}},

	// Invalid forms
	{``, false, nil, nil, []string{"mute"}},
	{`--2`, false, nil, nil, []string{}},
	{`+x`, false, nil, nil, []string{}},
	{`$1`, false, nil, nil, []string{}},
	{`mute more`, false, nil, nil, []string{}},
}

func TestVolume(t *testing.T) {
	commands.TestVerb(t, "volume", volumeTests)
}
