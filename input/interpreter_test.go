package input_test

import (
	"testing"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/options"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLISet tests that input.Interpreter registers a handler under the
// verb "set", dispatches the input line to this handler, and correctly
// manipulates the options table.
func TestCLISet(t *testing.T) {
	var err error

	a := api.NewTestAPI()
	opts := a.Options()
	iface := input.NewCLI(a)

	opts.Add(options.NewStringOption("foo"))
	err = opts.Get("foo").Set("this string must die")
	require.Nil(t, err)

	err = iface.Execute("set foo=something")
	assert.Nil(t, err)

	assert.Equal(t, "something", opts.Value("foo"))
}
