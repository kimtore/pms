package input_test

import (
	"testing"

	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/commands"
	"github.com/ambientsound/pms/options"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLISet tests that input.CLI registers a handler under the
// verb "set", dispatches the input line to this handler, and correctly
// manipulates the options table.
func TestCLISet(t *testing.T) {
	var err error

	api := commands.NewTestAPI()
	opts := api.Options()
	iface := input.NewCLI(api)

	iface.Register("set", commands.NewSet)

	opts.Add(options.NewStringOption("foo"))
	err = opts.Get("foo").Set("this string must die")
	require.Nil(t, err)

	err = iface.Execute("set foo=something")
	assert.Nil(t, err)

	assert.Equal(t, "something", opts.Value("foo"))
}
