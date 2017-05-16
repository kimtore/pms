package commands

import (
	"strings"
	"testing"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/stretchr/testify/assert"
)

// CommandTest is a structure for test data, and can be used to conveniently
// test Command instances.
type CommandTest struct {

	// The input data for the command, as seen on the command line.
	Input string

	// True if the command should parse and execute properly, false otherwise.
	Success bool

	// A callback function to call for every test, allowing customization of tests.
	Callback func(t *testing.T, cmd Command, api api.API, test CommandTest)

	// A slice of tab completion candidates to expect.
	TabComplete []string
}

// TestVerb runs table tests for Command implementations.
func TestVerb(t *testing.T, verb string, tests []CommandTest) {
	for n, test := range tests {
		api := api.NewTestAPI()
		cmd := New(verb, api)

		t.Logf("### Test %d: '%s'", n+1, test.Input)
		TestCommand(t, cmd, api, test)
	}
}

// TestCommand runs a single test a for Command implementation.
func TestCommand(t *testing.T, cmd Command, api api.API, test CommandTest) {
	reader := strings.NewReader(test.Input)
	scanner := lexer.NewScanner(reader)

	// Parse command
	cmd.SetScanner(scanner)
	err := cmd.Parse()

	// Test success
	if test.Success {
		assert.Nil(t, err, "Expected success when parsing '%s'", test.Input)
	} else {
		assert.NotNil(t, err, "Expected error when parsing '%s'", test.Input)
	}

	// Test tab completes
	completes := cmd.TabComplete()
	assert.Equal(t, test.TabComplete, completes)

	// Test callback function
	if test.Callback != nil {
		test.Callback(t, cmd, api, test)
	}
}
