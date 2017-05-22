package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var selectTests = []commands.Test{
	// Valid forms
	{`visual`, true, nil, nil, []string{}},
	{`toggle`, true, nil, nil, []string{}},
	{`nearby artist album`, true, nil, nil, []string{}},

	// Invalid forms
	{`foo`, false, nil, nil, []string{}},
	{`visual 1`, false, nil, nil, []string{}},
	{`toggle 1`, false, nil, nil, []string{}},
	{`nearby`, false, nil, nil, []string{}},

	// Tab completion
	{``, false, nil, nil, []string{
		"nearby",
		"toggle",
		"visual",
	}},
	{`t`, false, nil, nil, []string{
		"toggle",
	}},
}

func TestSelect(t *testing.T) {
	commands.TestVerb(t, "select", selectTests)
}
