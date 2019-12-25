package list_test

import (
	"github.com/ambientsound/pms/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumn(t *testing.T) {
	dataset := []string{
		"1",
		"12",
		"12345678",
		"123456789",
		"12345678901234567890123456789012345678901234567890123456789012345678901234567890",
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

	column.Add(string("foo"))
	assert.Equal(t, 3, column.Median())
	assert.Equal(t, 3, column.Avg())
	assert.Equal(t, 3, column.Max())

	column.Remove(string("foo"))
	assert.Equal(t, 0, column.Median())
	assert.Equal(t, 0, column.Avg())
	assert.Equal(t, 0, column.Max())

	column.Add(string("foo"))
	column.Add(string("foobarbaz"))
	assert.Equal(t, 6, column.Median())
	assert.Equal(t, 6, column.Avg())
	assert.Equal(t, 9, column.Max())

	column.Add(string("foo"))
	assert.Equal(t, 3, column.Median())
	assert.Equal(t, 5, column.Avg())
	assert.Equal(t, 9, column.Max())

	column.Remove(string("this item does not exist"))
	assert.Equal(t, 3, column.Median())
	assert.Equal(t, 5, column.Avg())
	assert.Equal(t, 9, column.Max())
}
