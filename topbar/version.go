package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/version"
)

// Version draws the short name of this application, as defined in the version module.
type Version struct {
	version string
}

func NewVersion(a api.API, param string) Fragment {
	return &Version{version.Version()}
}

func (w *Version) Text() (string, string) {
	return w.version, `version`
}
