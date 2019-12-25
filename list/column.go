package list

import (
	"sort"
)

type Item interface {
	Len() int
}

type Column struct {
	total   int
	sorted  bool
	lengths sort.IntSlice
}

func (c *Column) sort() {
	if !c.sorted {
		c.lengths.Sort()
		c.sorted = true
	}
}

func (c *Column) Add(item string) {
	ln := len(item)
	c.lengths = append(c.lengths, ln)
	c.total += len(item)
	c.sorted = false
}

func (c *Column) Avg() int {
	if c.lengths.Len() == 0 {
		return 0
	}
	return c.total / c.lengths.Len()
}

func (c *Column) Max() int {
	if len(c.lengths) == 0 {
		return 0
	}
	c.sort()
	return c.lengths[len(c.lengths)-1]
}

func (c *Column) Median() int {
	c.sort()
	ln := c.lengths.Len()
	mid := ln / 2
	if ln == 0 {
		return 0
	} else if ln%2 == 1 {
		return c.lengths[mid]
	}
	return (c.lengths[mid-1] + c.lengths[mid]) / 2
}

func (c *Column) Remove(item string) {
	idx := c.lengths.Search(len(item))
	if idx >= c.lengths.Len() {
		return
	} else if idx == c.lengths.Len()-1 {
		c.lengths = c.lengths[:idx]
	} else {
		c.lengths = append(c.lengths[:idx], c.lengths[idx+1:]...)
	}
	c.total -= len(item)
	c.sorted = false
}

func (c *Column) Set(items []string) {
	c.lengths = make(sort.IntSlice, len(items))
	c.sorted = false
	c.total = 0

	for i, item := range items {
		ln := len(item)
		c.lengths[i] = ln
		c.total += c.lengths[i]
	}

	c.sort()
}
