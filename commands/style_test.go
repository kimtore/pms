package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var styleTests = []commands.Test{
	// Valid forms
	{`stylekey`, true, nil, nil, []string{}},
	{`stylekey `, true, nil, nil, []string{"blink", "bold", "dim", "reverse", "underline"}},
	{`stylekey bar baz`, true, nil, nil, []string{}},
	{`stylekey color1 color2 blink bold dim reverse underline`, true, nil, nil, []string{"underline"}},
	{`stylekey blink color1 bold dim color2 reverse underline`, true, nil, nil, []string{"underline"}},

	// Invalid forms
	{``, false, nil, nil, []string{}},
	{`stylekey color1 color2 color3`, false, nil, nil, []string{}},
}

func TestStyle(t *testing.T) {
	commands.TestVerb(t, "style", styleTests)
}
