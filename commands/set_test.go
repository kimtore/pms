package commands_test

import (
	"testing"

	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/options"
	"github.com/stretchr/testify/assert"
)

var setTests = []commands.Test{
	// Valid forms
	{``, true, testSetInit, nil, []string{}},
	{`foo=bar`, true, testSetInit, testFooSet(`foo`, `bar`, true), []string{}},
	{`foo="bar baz"`, true, testSetInit, testFooSet(`foo`, `bar baz`, true), []string{}},
	{`foo=${}|;#`, true, testSetInit, testFooSet(`foo`, `${}|;`, true), []string{}},
	{`foo=x bar=x baz=x int=4 invbool`, true, testSetInit, testMultiSet, []string{}},
	{`foo=y foo`, true, testSetInit, testFooSet(`foo`, `y`, true), []string{}},

	// Invalid forms
	{`nonexist=foo`, true, testSetInit, testFooSet(`nonexist`, ``, false), []string{}},
	{`$=""`, false, testSetInit, nil, []string{}},
}

func TestSet(t *testing.T) {
	commands.TestVerb(t, "set", setTests)
}

func testSetInit(test *commands.TestData) {
	test.Api.Options().Add(options.NewStringOption("foo"))
	test.Api.Options().Add(options.NewStringOption("bar"))
	test.Api.Options().Add(options.NewStringOption("baz"))
	test.Api.Options().Add(options.NewIntOption("int"))
	test.Api.Options().Add(options.NewBoolOption("bool"))
}

func testFooSet(key, check string, ok bool) func(*commands.TestData) {
	return func(test *commands.TestData) {
		err := test.Cmd.Exec()
		assert.Equal(test.T, ok, err == nil, "Expected OK=%s", ok)
		if err != nil {
			return
		}
		val := test.Api.Options().StringValue(key)
		assert.Equal(test.T, check, val)
	}
}

func testMultiSet(test *commands.TestData) {
	err := test.Cmd.Exec()
	assert.Nil(test.T, err)
	opts := test.Api.Options()
	assert.Equal(test.T, "x", opts.StringValue("foo"))
	assert.Equal(test.T, "x", opts.StringValue("bar"))
	assert.Equal(test.T, "x", opts.StringValue("baz"))
	assert.Equal(test.T, 4, opts.IntValue("int"))
	assert.Equal(test.T, true, opts.BoolValue("bool"))
}
