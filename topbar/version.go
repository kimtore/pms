package topbar

import "github.com/ambientsound/pms/version"

// Version draws the current version tag.
type Version struct {
	fragment
}

func (w *Version) Width() int {
	return len(version.Version())
}

func (w *Version) Draw(x, y int) int {
	return w.drawNext(x, y, []rune(version.Version()), w.Style("version"))
}
