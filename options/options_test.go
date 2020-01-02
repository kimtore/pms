package options_test

import (
	"github.com/ambientsound/pms/options"
	"github.com/stretchr/testify/assert"
	"testing"
)

type printTest struct {
	key    string
	value  interface{}
	output string
}

var printTests = []printTest{
	{
		key:    `string`,
		value:  `bar`,
		output: `string="bar"`,
	},
	{
		key:    `int`,
		value:  56,
		output: `int=56`,
	},
	{
		key:    `bool`,
		value:  true,
		output: `bool`,
	},
	{
		key:    `bool`,
		value:  false,
		output: `nobool`,
	},
}

func TestPrint(t *testing.T) {
	for _, test := range printTests {
		output := options.Print(test.key, test.value)
		assert.Equal(t, test.output, output)
	}
}
