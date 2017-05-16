package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var selectTests = []commands.CommandTest{
	// Valid forms
	{`visual`, true, nil, []string{}},
	{`toggle`, true, nil, []string{}},

	// Invalid forms
	{`foo`, false, nil, []string{}},
	{`visual 1`, false, nil, []string{}},
	{`toggle 1`, false, nil, []string{}},

	// Tab completion
	{``, false, nil, []string{
		"toggle",
		"visual",
	}},
	{`t`, false, nil, []string{
		"toggle",
	}},
}

func TestSelect(t *testing.T) {
	commands.TestVerb(t, "select", selectTests)
}
