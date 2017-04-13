package input_test

import (
	"testing"

	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/commands"
	"github.com/ambientsound/pms/options"
	"github.com/stretchr/testify/assert"
)

// TestInterfaceSet tests that input.Interface registers a handler under the
// verb "set", dispatches the input line to this handler, and correctly
// manipulates the options table.
func TestInterfaceSet(t *testing.T) {
	var opt options.Option
	var err error

	opts := options.New()
	iface := input.NewInterface(&opts)

	iface.Register("set", commands.NewSet(&opts))

	if opt, err = options.NewStringOption("foo", "this string must die"); err != nil {
		t.Fatalf("Cannot add new string option: %s", err)
	}
	opts.Add(opt)

	err = iface.Execute("set foo=something")
	assert.Nil(t, err)

	assert.Equal(t, "something", opts.Value("foo"))
}
