package input_test

import (
	"github.com/spf13/viper"
	"testing"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input"
	"github.com/stretchr/testify/assert"
)

// TestCLISet tests that input.Interpreter registers a handler under the
// verb "set", dispatches the input line to this handler, and correctly
// manipulates the options table.
func TestCLISet(t *testing.T) {
	var err error

	a := api.NewTestAPI()
	opts := a.Options()
	iface := input.NewCLI(a)

	viper.Set("foo", "this string must die")

	err = iface.Exec("set foo=something")
	assert.Nil(t, err)

	assert.Equal(t, "something", opts.GetString("foo"))
}
