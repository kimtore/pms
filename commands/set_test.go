package commands_test

import (
	"testing"

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
	var token lexer.Token

	api := commands.NewTestAPI()
	opts := api.Options()

	opts.Add(options.NewStringOption("foo"))
	opts.Add(options.NewIntOption("intopt"))
	opts.Add(options.NewBoolOption("bar"))
	opts.Add(options.NewBoolOption("baz"))

	opts.Get("foo").Set("this string must die")
	opts.Get("intopt").Set("4")
	opts.Get("bar").Set("true")
	opts.Get("baz").Set("false")

	input_string := "foo=\"strings $are {}cool\" intopt=3 nobar invbaz"
	cmd := commands.NewSet(api)

	pos := 0
	npos := 0

	for {
		token, npos = lexer.NextToken(input_string[pos:])
		pos += npos
		err = cmd.Execute(token)
		if err != nil {
			t.Fatalf("Error while parsing input %s: %s", string(token.Runes), err)
		}
		if token.Class == lexer.TokenEnd {
			break
		}
	}

	assert.Equal(t, "strings $are {}cool", opts.Value("foo"))
	assert.Equal(t, 3, opts.Value("intopt"))
	assert.Equal(t, false, opts.Value("bar"))
	assert.Equal(t, true, opts.Value("baz"))
	assert.Equal(t, nil, opts.Value("skrot"))
}
