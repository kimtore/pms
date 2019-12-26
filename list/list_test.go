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

func TestListNextOf(t *testing.T) {
	lst := list.New()

	dataset := []list.Row{
		{
			"foo": "x",
			"bar": "x",
		},
		{
			"foo": "x",
			"bar": "xyz",
		},
		{
			"foo": "foo",
			"bar": "x",
		},
		{
			"foo": "x",
			"bar": "x",
		},
		{
			"foo": "x",
			"bar": "x",
		},
	}

	for _, row := range dataset {
		lst.Add(row)
	}

	t.Run("stop in the middle at positive direction", func(t *testing.T) {
		next := lst.NextOf([]string{"foo"}, 0, 1)
		assert.Equal(t, 2, next)

		next = lst.NextOf([]string{"foo"}, 1, 1)
		assert.Equal(t, 2, next)
	})

	t.Run("stop in the middle at negative direction", func(t *testing.T) {
		next := lst.NextOf([]string{"foo"}, 4, -1)
		assert.Equal(t, 3, next)
	})

	t.Run("stop at the first multi-tag diff", func(t *testing.T) {
		next := lst.NextOf([]string{"foo", "bar"}, 0, 1)
		assert.Equal(t, 1, next)
	})

	t.Run("stop at the end if no diff found", func(t *testing.T) {
		next := lst.NextOf([]string{"bar"}, 2, 1)
		assert.Equal(t, 5, next)
	})
}
