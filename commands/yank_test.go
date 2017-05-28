package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var yankTests = []commands.Test{
	// Yank takes to parameters.
	{``, true, nil, nil, []string{}},
	{`    `, true, nil, nil, []string{}},

	// Invalid forms
	{`foo`, false, nil, nil, []string{}},
	{`foo bar`, false, nil, nil, []string{}},
}

func TestYank(t *testing.T) {
	commands.TestVerb(t, "yank", yankTests)
}
