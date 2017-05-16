package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var selectTests = []commands.CommandTest{
	// Valid forms
	{`visual`, true, []string{}},
	{`toggle`, true, []string{}},

	// Invalid forms
	{`foo`, false, []string{}},
	{`visual 1`, false, []string{}},
	{`toggle 1`, false, []string{}},

	// Tab completion
	{``, false, []string{
		"toggle",
		"visual",
	}},
	{`t`, false, []string{
		"toggle",
	}},
}

func TestSelect(t *testing.T) {
	commands.TestVerb(t, "select", selectTests)
}
