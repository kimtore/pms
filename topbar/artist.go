package topbar

import (
	"github.com/ambientsound/pms/api"
)

// Artist draws the artist tag.
type Artist struct {
	tag
}

func NewArtist(a api.API) Fragment {
	return &Artist{
		tag{"artist", fragment{api: a}},
	}
}
