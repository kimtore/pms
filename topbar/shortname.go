package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/version"
)

// Shortname draws the short name of this application, as defined in the version module.
type Shortname struct {
	fragment
}

func NewShortname(a api.API) Fragment {
	return &Shortname{
		fragment{api: a},
	}
}

func (w *Shortname) Width() int {
	return len(version.ShortName())
}

func (w *Shortname) Draw(x, y int) int {
	return w.drawNext(x, y, []rune(version.ShortName()), w.Style("shortName"))
}
