package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var isolateTests = []commands.Test{
	// Valid forms
	{`artist`, true, nil, nil, []string{}},
	{`artist t`, true, initSongTags, nil, []string{"title"}},
	{`artist tr$ack title`, true, initSongTags, nil, []string{"title"}},

	// Invalid forms
	{``, false, nil, nil, []string{}},
}

func TestIsolate(t *testing.T) {
	commands.TestVerb(t, "isolate", isolateTests)
}
