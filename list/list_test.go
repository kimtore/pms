package list_test

import (
	"github.com/ambientsound/pms/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {
	lst := list.New()

	dataset := []list.Row{
		{
			"foo": "foo content",
			"bar": "nope",
		},
		{
			"baz": "foo",
			"bar": "nopenope",
		},
		{
			"bar": "hope",
		},
	}

	for _, row := range dataset {
		lst.Add(row)
	}

	t.Run("dataset content is correct", func(t *testing.T) {
		assert.Equal(t, len(dataset), lst.Len())

		for i := 0; i < lst.Len(); i++ {
			assert.Equal(t, dataset[i], lst.Row(i))
		}

	})

	t.Run("column names are correct", func(t *testing.T) {
		names := lst.ColumnNames()
		assert.ElementsMatch(t, []string{"bar", "baz", "foo"}, names)
	})

	t.Run("columns have correct content and sizes", func(t *testing.T) {
		names := []string{"foo", "bar", "baz"}
		cols := lst.Columns(names)

		assert.Len(t, cols, len(names))

		assert.Equal(t, 11, cols[0].Median())
		assert.Equal(t, 11, cols[0].Avg())

		assert.Equal(t, 4, cols[1].Median())
		assert.Equal(t, 5, cols[1].Avg())

		assert.Equal(t, 3, cols[2].Median())
		assert.Equal(t, 3, cols[2].Avg())
	})
}
