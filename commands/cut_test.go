package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var cutTests = []commands.Test{
	// Cut takes to parameters.
	{``, true, nil, nil, []string{}},
	{`    `, true, nil, nil, []string{}},

	// Invalid forms
	{`foo`, false, nil, nil, []string{}},
	{`foo bar`, false, nil, nil, []string{}},
}

func TestCut(t *testing.T) {
	commands.TestVerb(t, "cut", cutTests)
}
