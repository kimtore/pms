package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var removeTests = []commands.Test{
	// Remove takes to parameters.
	{``, true, nil, nil, []string{}},
	{`    `, true, nil, nil, []string{}},

	// Invalid forms
	{`foo`, false, nil, nil, []string{}},
	{`foo bar`, false, nil, nil, []string{}},
}

func TestRemove(t *testing.T) {
	commands.TestVerb(t, "remove", removeTests)
}
