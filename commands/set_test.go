package commands_test

import (
	"strings"
	"testing"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/options"
	"github.com/stretchr/testify/assert"
)

// TestSet tests the pms.input.commands.Set.Execute() function. A string of input
// parameters are given, and Execute() is expected to populate a options.Options
// struct with parsed values.
func TestSet(t *testing.T) {
	var err error

	a := api.NewTestAPI()
	opts := a.Options()

	opts.Add(options.NewStringOption("foo"))
	opts.Add(options.NewIntOption("intopt"))
	opts.Add(options.NewBoolOption("bar"))
	opts.Add(options.NewBoolOption("baz"))

	opts.Get("foo").Set("this string must die")
	opts.Get("intopt").Set("4")
	opts.Get("bar").Set("true")
	opts.Get("baz").Set("false")

	line := `foo="strings $are {}cool" intopt=3 nobar invbaz`
	cmd := commands.NewSet(a)

	reader := strings.NewReader(line)
	scanner := lexer.NewScanner(reader)

	for {
		class, token := scanner.Scan()

		err = cmd.Execute(class, token)

		assert.Nil(t, err, "Error while parsing input '%s' of class %d: %s", token, class, err)

		if class == lexer.TokenEnd {
			break
		}
	}

	assert.Equal(t, "strings $are {}cool", opts.Value("foo"))
	assert.Equal(t, 3, opts.Value("intopt"))
	assert.Equal(t, false, opts.Value("bar"))
	assert.Equal(t, true, opts.Value("baz"))
	assert.Equal(t, nil, opts.Value("skrot"))
}
