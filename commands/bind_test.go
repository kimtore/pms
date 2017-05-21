package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var bindTests = []commands.Test{
	// Valid forms
	{`foo bar`, true, nil, nil, []string{}},
	{`foo bar baz`, true, nil, nil, []string{}},
	{`[]{}$|"test" foo bar`, true, nil, nil, []string{}},

	// Invalid forms
	{``, false, nil, nil, []string{}},
	{`x`, false, nil, nil, []string{}},
}

func TestBind(t *testing.T) {
	commands.TestVerb(t, "bind", bindTests)
}
