package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var seekTests = []commands.Test{
	// Valid forms
	{`-2`, true, nil, nil, []string{}},
	{`+13`, true, nil, nil, []string{}},
	{`1329`, true, nil, nil, []string{}},

	// Invalid forms
	{`nan`, false, nil, nil, []string{}},
	{`+++1`, false, nil, nil, []string{}},
	{`-foo`, false, nil, nil, []string{}},
	{`$1`, false, nil, nil, []string{}},
}

func TestSeek(t *testing.T) {
	commands.TestVerb(t, "seek", seekTests)
}
