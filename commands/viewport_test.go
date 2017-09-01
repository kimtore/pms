package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var viewportTests = []commands.Test{
	// Valid forms
	{`up`, true, nil, nil, []string{}},
	{`down`, true, nil, nil, []string{}},

	// Invalid forms
	{`up 1`, false, nil, nil, []string{}},
	{`down 1`, false, nil, nil, []string{}},
	{`nonsense`, false, nil, nil, []string{}},

	// Tab completion
	{``, false, nil, nil, []string{
		"down",
		"up",
	}},
	{`u`, false, nil, nil, []string{
		"up",
	}},
	{`do`, false, nil, nil, []string{
		"down",
	}},
}

func TestViewport(t *testing.T) {
	commands.TestVerb(t, "viewport", viewportTests)
}
