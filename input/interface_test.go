package input_test

import (
	"testing"

	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/commands"
	"github.com/ambientsound/pms/options"
	"github.com/stretchr/testify/assert"
)

// TestCLISet tests that input.CLI registers a handler under the
// verb "set", dispatches the input line to this handler, and correctly
// manipulates the options table.
func TestCLISet(t *testing.T) {
	var err error

	opts := options.New()
	iface := input.NewCLI()
	messages := make(chan string, 1024)

	iface.Register("set", commands.NewSet(opts, messages))

	opts.Add(options.NewStringOption("foo", "this string must die"))

	err = iface.Execute("set foo=something")
	assert.Nil(t, err)

	assert.Equal(t, "something", opts.Value("foo"))
}
