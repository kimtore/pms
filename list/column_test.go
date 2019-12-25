package list_test

import (
	"github.com/ambientsound/pms/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

type stringItem string

func (s stringItem) Len() int {
	return len(s)
}

func TestColumn(t *testing.T) {
	dataset := []list.Item{
		stringItem("1"),
		stringItem("12"),
		stringItem("12345678"),
		stringItem("123456789"),
		stringItem("12345678901234567890123456789012345678901234567890123456789012345678901234567890"),
	}

	column := &list.Column{}

	assert.Equal(t, 0, column.Median())
	assert.Equal(t, 0, column.Avg())

	column.Set(dataset)

	assert.Equal(t, 20, column.Avg())
	assert.Equal(t, 8, column.Median())
}

func TestColumnAddRemove(t *testing.T) {
	column := &list.Column{}

	column.Add(stringItem("foo"))
	assert.Equal(t, 3, column.Median())
	assert.Equal(t, 3, column.Avg())

	column.Remove(stringItem("foo"))
	assert.Equal(t, 0, column.Median())
	assert.Equal(t, 0, column.Avg())

	column.Add(stringItem("foo"))
	column.Add(stringItem("foobarbaz"))
	assert.Equal(t, 6, column.Median())
	assert.Equal(t, 6, column.Avg())

	column.Add(stringItem("foo"))
	assert.Equal(t, 3, column.Median())
	assert.Equal(t, 5, column.Avg())

	column.Remove(stringItem("this item does not exist"))
	assert.Equal(t, 3, column.Median())
	assert.Equal(t, 5, column.Avg())
}
