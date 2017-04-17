package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/input/commands"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/options"
	"github.com/stretchr/testify/assert"
)

// TestSet tests the pms.input.commands.Set.Execute() function. A string of input
// parameters are given, and Execute() is expected to populate a options.Options
// struct with parsed values.
func TestSet(t *testing.T) {
	var opt options.Option
	var err error
	var token lexer.Token

	opts := options.New()

	if opt, err = options.NewStringOption("foo", "this string must die"); err != nil {
		t.Fatalf("Cannot add new string option: %s", err)
	}
	opts.Add(opt)

	if opt, err = options.NewIntOption("intopt", 4); err != nil {
		t.Fatalf("Cannot add new integer option: %s", err)
	}
	opts.Add(opt)

	if opt, err = options.NewBoolOption("bar", true); err != nil {
		t.Fatalf("Cannot add new boolean option: %s", err)
	}
	opts.Add(opt)

	if opt, err = options.NewBoolOption("baz", false); err != nil {
		t.Fatalf("Cannot add new boolean option: %s", err)
	}
	opts.Add(opt)

	input_string := "foo=bar intopt=3 nobar invbaz"
	cmd := commands.NewSet(opts)

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

	assert.Equal(t, "bar", opts.Value("foo"))
	assert.Equal(t, 3, opts.Value("intopt"))
	assert.Equal(t, false, opts.Value("bar"))
	assert.Equal(t, true, opts.Value("baz"))
	assert.Equal(t, nil, opts.Value("skrot"))
}
