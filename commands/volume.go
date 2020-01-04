package commands

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/input/lexer"
)

var preMuteVolume int

// Volume adjusts MPD's volume.
type Volume struct {
	command
	api      api.API
	sign     int
	volume   int
	finished bool
	mute     bool
}

// NewVolume returns Volume.
func NewVolume(api api.API) Command {
	return &Volume{
		api: api,
	}
}

// Parse implements Command.
func (cmd *Volume) Parse() error {

	playerStatus := cmd.api.PlayerStatus()

	tok, lit := cmd.ScanIgnoreWhitespace()
	cmd.setTabComplete(lit, []string{"mute"})

	// Check for muted status.
	if tok == lexer.TokenIdentifier && lit == "mute" {
		cmd.mute = true
		cmd.setTabCompleteEmpty()
		return cmd.ParseEnd()
	}

	// If not muted, try to parse a number.
	cmd.Unscan()
	_, ilit, absolute, err := cmd.ParseInt()
	if err != nil {
		return err
	}

	if absolute {
		cmd.volume = ilit
	} else {
		cmd.volume = int(playerStatus.Device.Volume) + ilit
	}

	cmd.validateVolume()

	cmd.setTabCompleteEmpty()
	return cmd.ParseEnd()
}

// validateVolume clamps the volume to the allowable range
func (cmd *Volume) validateVolume() {
	if cmd.volume > 100 {
		cmd.volume = 100
	} else if cmd.volume < 0 {
		cmd.volume = 0
	}
}

// Exec implements Command.
func (cmd *Volume) Exec() error {
	client, err := cmd.api.Spotify()
	if err != nil {
		return err
	}

	playerStatus := cmd.api.PlayerStatus()

	switch {
	case cmd.mute && playerStatus.Device.Volume == 0:
		cmd.volume = preMuteVolume
	case cmd.mute && playerStatus.Device.Volume > 0:
		preMuteVolume = playerStatus.Device.Volume
		cmd.volume = 0
	}

	return client.Volume(cmd.volume)
}
