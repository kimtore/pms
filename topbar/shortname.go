package topbar

import "github.com/ambientsound/pms/version"

// Shortname draws the short name of this application, as defined in the version module.
type Shortname struct {
	fragment
}

func NewShortname() Fragment {
	return &Shortname{
	//fragment{api: api},
	}
}

func (w *Shortname) Width() int {
	return len(version.ShortName())
}

func (w *Shortname) Draw(x, y int) int {
	return w.drawNext(x, y, []rune(version.ShortName()), w.Style("shortName"))
}
