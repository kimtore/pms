package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var unbindTests = []commands.Test{
	// Valid forms
	{`f`, true, nil, nil, []string{}},
	{`foo`, true, nil, nil, []string{}},
	{`[]{}$|"test"`, true, nil, nil, []string{}},

	// Invalid forms
	{``, false, nil, nil, []string{}},
	{`foo bar`, false, nil, nil, []string{}},
	{`foo bar baz`, false, nil, nil, []string{}},
}

func TestUnbind(t *testing.T) {
	commands.TestVerb(t, "unbind", unbindTests)
}
