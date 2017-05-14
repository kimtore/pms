package commands_test

import (
	"strings"
	"testing"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/stretchr/testify/assert"
)

var cursorTests = []struct {
	input       string
	success     bool
	tabComplete []string
}{
	// Valid forms
	{`up`, true, []string{}},
	{`down`, true, []string{}},
	//{`pgup`, true},
	//{`pageup`, true},
	//{`pagedn`, true},
	//{`pagedown`, true},
	{`home`, true, []string{}},
	{`end`, true, []string{}},
	{`current`, true, []string{}},
	{`random`, true, []string{}},
	{`next-of tag1,tag2`, true, []string{}},
	{`prev-of tag1,tag2`, true, []string{}},

	// Invalid forms
	{`up 1`, false, []string{}},
	{`down 1`, false, []string{}},
	//{`pgup 1`, false},
	//{`pageup 1`, false},
	//{`pagedn 1`, false},
	//{`pagedown 1`, false},
	{`home 1`, false, []string{}},
	{`end 1`, false, []string{}},
	{`current 1`, false, []string{}},
	{`random 1`, false, []string{}},
	{`next-of`, false, []string{"artist", "title"}},
	{`prev-of`, false, []string{"artist", "title"}},
	{`next-of 1 2`, false, []string{}},
	{`prev-of 1 2`, false, []string{}},

	// Tab completion
	{``, false, []string{
		"current",
		"down",
		"end",
		"home",
		"next-of",
		"pagedn",
		"pagedown",
		"pageup",
		"pgdn",
		"pgup",
		"prev-of",
		"random",
		"up",
	}},
	{`page`, false, []string{
		"pagedn",
		"pagedown",
		"pageup",
	}},
}

func TestCursor(t *testing.T) {
	for n, test := range cursorTests {

		api := api.NewTestAPI()
		cmd := commands.NewCursor(api)

		t.Logf("### Test %d: '%s'", n+1, test.input)

		reader := strings.NewReader(test.input)
		scanner := lexer.NewScanner(reader)

		// Parse command
		err := cmd.Parse(scanner)

		// Test success
		if test.success {
			assert.Nil(t, err, "Expected success when parsing '%s'", test.input)
		} else {
			assert.NotNil(t, err, "Expected error when parsing '%s'", test.input)
		}

		// Test tab completes
		completes := cmd.TabComplete()
		assert.Equal(t, test.tabComplete, completes)
	}
}
