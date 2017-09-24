package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var styleTests = []commands.Test{
	// Valid forms
	{`stylekey`, true, nil, nil, []string{}},
	{`stylekey `, true, nil, nil, []string{"bold", "reverse", "underline"}},
	{`stylekey bar baz`, true, nil, nil, []string{}},
	{`stylekey color1 color2 bold reverse underline`, true, nil, nil, []string{"underline"}},
	{`stylekey color1 bold color2 reverse underline`, true, nil, nil, []string{"underline"}},

	// Invalid forms
	{``, false, nil, nil, []string{}},
	{`stylekey color1 color2 color3`, false, nil, nil, []string{}},
}

func TestStyle(t *testing.T) {
	commands.TestVerb(t, "style", styleTests)
}
