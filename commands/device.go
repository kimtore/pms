package commands

import (
	"fmt"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/spotify/devices"
	"github.com/zmb3/spotify"
)

// Device seeks forwards or backwards in the currently playing track.
type Device struct {
	command
	api        api.API
	deviceID   spotify.ID
	deviceName string
}

// NewDevice returns Device.
func NewDevice(api api.API) Command {
	return &Device{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Device) Parse() error {
	tok, lit := cmd.ScanIgnoreWhitespace()

	switch tok {
	case lexer.TokenIdentifier:
	default:
		return fmt.Errorf("unexpected '%s', expected identifier", lit)
	}

	if lit != "activate" {
		return fmt.Errorf("unexpected '%s', expected 'activate'", lit)
	}

	tok, lit = cmd.ScanIgnoreWhitespace()

	switch tok {
	case lexer.TokenEnd:
		lst, ok := cmd.api.List().(*spotify_devices.List)
		if !ok {
			return fmt.Errorf("must be run in the devices window unless device ID is specified")
		}
		device := lst.CursorDevice()
		if device == nil {
			return fmt.Errorf("no devices available")
		}

		cmd.deviceID = device.ID
		cmd.deviceName = device.Name
		return nil

	case lexer.TokenIdentifier:
		cmd.deviceID = spotify.ID(lit)
		cmd.deviceName = "device ID " + lit
	}

	return cmd.ParseEnd()
}

// Exec implements Command.
func (cmd *Device) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	log.Infof("Transferring playback to %s...", cmd.deviceName)

	return client.TransferPlayback(cmd.deviceID, cmd.api.PlayerStatus().Playing)
}
