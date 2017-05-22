package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var playTests = []commands.Test{
	// Valid forms
	{``, true, nil, nil, []string{"cursor", "selection"}},
	{`cursor`, true, nil, nil, []string{}},
	{`selection`, true, nil, nil, []string{}},

	// Invalid forms
	{`foo`, false, nil, nil, []string{}},
	{`cursor 1`, false, nil, nil, []string{}},
	{`selection 1`, false, nil, nil, []string{}},
}

func TestPlay(t *testing.T) {
	commands.TestVerb(t, "play", playTests)
}
