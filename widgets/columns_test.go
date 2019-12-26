package widgets_test

import (
	"github.com/ambientsound/pms/widgets"
	"github.com/stretchr/testify/assert"
	"testing"
)

type columnTitleTest struct {
	input  string
	output string
}

var columnTitleTests = []columnTitleTest{
	{"", ""},
	{"foo", "Foo"},
	{"fooBar", "Foo Bar"},
	{"fooBarBaz", "Foo Bar Baz"},
}

func TestColumnTitle(t *testing.T) {
	for _, test := range columnTitleTests {
		output := widgets.ColumnTitle(test.input)
		assert.Equal(t, test.output, output)
	}
}
