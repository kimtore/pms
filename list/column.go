package list

import (
	"sort"
)

type Column struct {
	name    string
	total   int
	sorted  bool
	lengths sort.IntSlice
}

func (c *Column) Add(item Item) {
	c.lengths = append(c.lengths, item.Len())
	c.total += item.Len()
	c.sorted = false
	c.lengths.Sort()
}

func (c *Column) Avg() int {
	if c.lengths.Len() == 0 {
		return 0
	}
	return c.total / c.lengths.Len()
}

func (c *Column) Median() int {
	if !c.sorted {
		c.lengths.Sort()
		c.sorted = true
	}
	ln := c.lengths.Len()
	mid := ln / 2
	if ln == 0 {
		return 0
	} else if ln%2 == 1 {
		return c.lengths[mid]
	}
	return (c.lengths[mid-1] + c.lengths[mid]) / 2
}

func (c *Column) Name() string {
	return c.name
}

func (c *Column) Remove(item Item) {
	idx := c.lengths.Search(item.Len())
	if idx >= c.lengths.Len() {
		return
	} else if idx == c.lengths.Len()-1 {
		c.lengths = c.lengths[:idx]
	} else {
		c.lengths = append(c.lengths[:idx], c.lengths[idx+1:]...)
	}
	c.total -= item.Len()
	c.sorted = false
}

func (c *Column) Set(items []Item) {
	c.lengths = make(sort.IntSlice, len(items))
	c.total = 0

	for i, item := range items {
		c.lengths[i] = item.Len()
		c.total += c.lengths[i]
	}

	c.lengths.Sort()
	c.sorted = true
}

func (c *Column) SetName(name string) {
	c.name = name
}
