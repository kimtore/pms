package topbar

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/version"
)

// Version draws the current version tag.
type Version struct {
	fragment
}

func NewVersion(a api.API, param string) Fragment {
	return &Version{
		fragment{api: a},
	}
}

func (w *Version) Width() int {
	return len(version.Version())
}

func (w *Version) Draw(x, y int) int {
	return w.drawNextString(x, y, version.Version(), w.Style("version"))
}
