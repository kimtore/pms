package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var singleTests = []commands.Test{
	// Valid forms
	{`on`, true, nil, nil, []string{}},
	{`off`, true, nil, nil, []string{}},
	{`toggle`, true, nil, nil, []string{}},

	// Invalid forms
	{`--2`, false, nil, nil, []string{}},
	{`+x`, false, nil, nil, []string{}},
	{`$1`, false, nil, nil, []string{}},
	{`on off`, false, nil, nil, []string{}},

	// Tab completion
	{``, true, nil, nil, []string{
		"on",
		"off",
		"toggle",
	}},
	{`t`, false, nil, nil, []string{
		"toggle",
	}},
	{`o`, false, nil, nil, []string{
		"on",
		"off",
	}},
}

func TestSingle(t *testing.T) {
	commands.TestVerb(t, "single", singleTests)
}
