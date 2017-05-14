package commands

import (
	"fmt"
	"strconv"

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

func NewVolume(api api.API) Command {
	return &Volume{
		api: api,
	}
}

func (cmd *Volume) Execute(class int, s string) error {
	var err error

	switch class {
	case lexer.TokenIdentifier:

		if cmd.finished {
			return fmt.Errorf("Unexpected '%s', expected END", s)
		}

		switch {
		case s == "mute":
			cmd.mute = true
			cmd.finished = true
			return nil
		case s[0] == '+':
			cmd.sign = 1
			s = s[1:]
		case s[0] == '-':
			cmd.sign = -1
			s = s[1:]
		}

		cmd.volume, err = strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Unexpected '%s', expected number", s)
		}

		cmd.finished = true

	case lexer.TokenEnd:
		if !cmd.finished {
			return fmt.Errorf("Unexpected END, expected absolute or relative volume")
		}

		client := cmd.api.MpdClient()
		if client == nil {
			return fmt.Errorf("Unable to control volume: cannot communicate with MPD")
		}
		status := cmd.api.PlayerStatus()

		switch {
		case cmd.mute && status.Volume == 0:
			cmd.volume = preMuteVolume
		case cmd.mute && status.Volume > 0:
			preMuteVolume = status.Volume
			cmd.volume = 0
		case cmd.sign != 0:
			cmd.volume *= cmd.sign
			cmd.volume = status.Volume + cmd.volume
		}

		return client.SetVolume(cmd.volume)

	default:
		return fmt.Errorf("Unknown input '%s', expected END", s)
	}

	return nil
}
