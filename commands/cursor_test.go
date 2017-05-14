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
	input   string
	success bool
}{
	// Valid forms
	{`up`, true},
	{`down`, true},
	//{`pgup`, true},
	//{`pageup`, true},
	//{`pagedn`, true},
	//{`pagedown`, true},
	{`home`, true},
	{`end`, true},
	{`current`, true},
	{`random`, true},
	{`next-of tag1,tag2`, true},
	{`prev-of tag1,tag2`, true},

	// Invalid forms
	{`up 1`, false},
	{`down 1`, false},
	//{`pgup 1`, false},
	//{`pageup 1`, false},
	//{`pagedn 1`, false},
	//{`pagedown 1`, false},
	{`home 1`, false},
	{`end 1`, false},
	{`current 1`, false},
	{`random 1`, false},
	{`next-of`, false},
	{`prev-of`, false},
	{`next-of 1 2`, false},
	{`prev-of 1 2`, false},
}

func TestCursor(t *testing.T) {
	for n, test := range cursorTests {

		api := api.NewTestAPI()
		cmd := commands.NewCursor(api)

		t.Logf("### Test %d: '%s'", n+1, test.input)

		reader := strings.NewReader(test.input)
		scanner := lexer.NewScanner(reader)
		err := cmd.Parse(scanner)

		if test.success {
			assert.Nil(t, err, "Expected success when parsing '%s'", test.input)
		} else {
			assert.NotNil(t, err, "Expected error when parsing '%s'", test.input)
		}
	}
}
