package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var addTests = []commands.CommandTest{
	// Valid forms
	{``, true, []string{}},
	{`foo bar baz`, true, []string{}},
	{`http://example.com/stream.mp3?foo=bar&baz=foo foo bar baz`, true, []string{}},
	{`|`, true, []string{}},
	{`|{}$`, true, []string{}},

	// No invalid forms, all input is accepted
}

func TestAdd(t *testing.T) {
	commands.TestVerb(t, "add", addTests)
}
