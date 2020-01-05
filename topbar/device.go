package topbar

import (
	"strings"

	"github.com/ambientsound/pms/api"
)

// Device shows information about the currently playing music device.
type Device struct {
	api   api.API
	param string
	style string
}

// NewDevice returns Device.
func NewDevice(a api.API, param string) Fragment {
	if len(param) == 0 {
		param = "name"
	}
	style := `device` + strings.Title(param)
	return &Device{a, param, style}
}

// Text implements Fragment.
func (w *Device) Text() (string, string) {
	dev := w.api.PlayerStatus().Device
	switch w.param {
	default:
		fallthrough
	case "name":
		return dev.Name, w.style
	case "type":
		return dev.Type, w.style
	}
}
