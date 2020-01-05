package db

import (
	"github.com/ambientsound/pms/list"
	"strconv"
)

type List struct {
	list.Base
	lists map[string]list.List
}

var _ list.List = &List{}

func New() *List {
	this := &List{}
	this.Clear()
	this.SetID("windows")
	this.SetName("Windows")
	this.SetVisibleColumns([]string{"name", "size"})
	this.lists = make(map[string]list.List)
	return this
}

func Row(lst list.List) list.Row {
	return list.Row{
		list.RowIDKey: lst.ID(),
		"name":        lst.Name(),
		"size":        strconv.Itoa(lst.Len()),
	}
}

// Cache adds a list to the database. Returns the row number of the list.
func (s *List) Cache(lst list.List) int {
	existing, ok := s.lists[lst.ID()]
	s.lists[lst.ID()] = lst
	if !ok {
		s.Add(Row(lst))
		return s.Len() - 1
	}
	rown, _ := s.RowNum(existing.ID())
	row := s.Row(rown)
	row["size"] = Row(existing)["size"]
	return rown
}

func (s *List) Current() list.List {
	row := s.CursorRow()
	if row == nil {
		return nil
	}
	return s.lists[row.ID()]
}

func (s *List) List(id string) list.List {
	return s.lists[id]
}
