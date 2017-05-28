package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
)

var pasteTests = []commands.Test{
	// Valid forms
	{``, true, nil, nil, []string{"after", "before"}},
	{`before`, true, nil, nil, []string{}},
	{`after`, true, initSongTags, nil, []string{}},

	// Invalid forms
	{`bef`, false, nil, nil, []string{"before"}},
	{`before the apocalypse`, false, nil, nil, []string{}},
	{`after midnight`, false, nil, nil, []string{}},
}

func TestPaste(t *testing.T) {
	commands.TestVerb(t, "paste", pasteTests)
}
