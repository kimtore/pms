package commands_test

import (
	"strings"
	"testing"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/stretchr/testify/assert"
)

var selectTests = []struct {
	input       string
	success     bool
	tabComplete []string
}{
	// Valid forms
	{`visual`, true, []string{}},
	{`toggle`, true, []string{}},

	// Invalid forms
	{`foo`, false, []string{}},
	{`visual 1`, false, []string{}},
	{`toggle 1`, false, []string{}},

	// Tab completion
	{``, false, []string{
		"toggle",
		"visual",
	}},
	{`t`, false, []string{
		"toggle",
	}},
}

func TestSelect(t *testing.T) {
	for n, test := range selectTests {

		api := api.NewTestAPI()
		cmd := commands.NewSelect(api)

		t.Logf("### Test %d: '%s'", n+1, test.input)

		reader := strings.NewReader(test.input)
		scanner := lexer.NewScanner(reader)

		// Parse command
		cmd.SetScanner(scanner)
		err := cmd.Parse()

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
