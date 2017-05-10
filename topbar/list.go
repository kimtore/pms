package topbar

import (
	"fmt"

	"github.com/ambientsound/pms/api"
)

// List draws information about the current songlist.
type List struct {
	api api.API
	f   func() (string, string)
}

// NewList returns List.
func NewList(a api.API, param string) Fragment {
	list := &List{a, nil}
	switch param {
	case `index`:
		list.f = list.textIndex
	case `title`:
		list.f = list.textTitle
	case `total`:
		list.f = list.textTotal
	default:
		list.f = list.textNone
	}
	return list
}

// Text implements Fragment.
func (w *List) Text() (string, string) {
	return w.f()
}

func (w *List) textNone() (string, string) {
	return ``, ``
}

func (w *List) textIndex() (string, string) {
	index, err := w.api.SonglistWidget().SonglistIndex()
	if err == nil {
		return fmt.Sprintf("%d", index+1), `listIndex`
	} else {
		return `new`, `listIndex`
	}
}

func (w *List) textTotal() (string, string) {
	total := w.api.SonglistWidget().SonglistsLen()
	return fmt.Sprintf("%d", total), `listTotal`
}

func (w *List) textTitle() (string, string) {
	return w.api.SonglistWidget().Songlist().Name(), `listTitle`
}
