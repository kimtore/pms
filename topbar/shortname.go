package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/version"
)

// Shortname draws the short name of this application, as defined in the version module.
type Shortname struct {
	shortname string
}

func NewShortname(a api.API, param string) Fragment {
	return &Shortname{version.ShortName()}
}

func (w *Shortname) Text() (string, string) {
	return w.shortname, `shortName`
}
