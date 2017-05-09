package topbar_test

import (
	"testing"

	"github.com/ambientsound/pms/topbar"
	"github.com/stretchr/testify/assert"
)

type result struct {
	class int
	str   string
}

var topbarTests = []struct {
	input   string
	success bool
	width   int
	height  int
}{
	{"plain", true, 1, 1},
	{"plain|white", true, 2, 1},
	{"plain|white|tests;multiple lines", true, 3, 2},
	{";;more;lines|here", true, 2, 4},
	{"$shortname|$version", true, 2, 1},
	{"$bogus_variable", false, 0, 0},
	{"$$var1", false, 0, 0},
	{"{}", false, 0, 0},
	{"# comment", false, 0, 0},
	{"\"quoted\"", true, 1, 1},
}

func TestTopbarCount(t *testing.T) {
	for _, test := range topbarTests {

		matrix, err := topbar.Parse(test.input)
		if test.success {
			assert.Nil(t, err, "Expected success in topbar parser when parsing '%s'", test.input)
		} else {
			assert.NotNil(t, err, "Expected error in topbar parser when parsing '%s'", test.input)
			continue
		}

		assert.Equal(t, test.height, len(matrix),
			"Topbar input '%s' should yield %d lines, got %d instead", test.input, test.height, len(matrix))

		for y := 0; y < len(matrix); y++ {
			assert.Equal(t, test.width, len(matrix[y]),
				"Topbar input '%s' should yield %d columns on line %d, got %d instead", test.input, test.width, y+1, len(matrix[y]))
		}
	}
}
