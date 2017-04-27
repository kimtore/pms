package commands

import (
	"fmt"
	"strconv"

	"github.com/ambientsound/pms/input/lexer"

	"github.com/ambientsound/gompd/mpd"
	pms_mpd "github.com/ambientsound/pms/mpd"
)

// Volume adjusts MPD's volume.
type Volume struct {
	mpdClient     func() *mpd.Client
	mpdStatus     func() pms_mpd.PlayerStatus
	sign          int
	volume        int
	finished      bool
	mute          bool
	preMuteVolume int
}

func NewVolume(mpdClient func() *mpd.Client, mpdStatus func() pms_mpd.PlayerStatus) *Volume {
	return &Volume{mpdClient: mpdClient, mpdStatus: mpdStatus}
}

func (cmd *Volume) Reset() {
	cmd.sign = 0
	cmd.volume = 0
	cmd.finished = false
	cmd.mute = false
	// DO NOT reset preMuteVolume
}

func (cmd *Volume) Execute(t lexer.Token) error {
	var err error

	switch t.Class {
	case lexer.TokenIdentifier:
		s := t.String()

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

		client := cmd.mpdClient()
		if client == nil {
			return fmt.Errorf("Unable to control volume: cannot communicate with MPD")
		}
		status := cmd.mpdStatus()

		switch {
		case cmd.mute && status.Volume == 0:
			cmd.volume = cmd.preMuteVolume
		case cmd.mute && status.Volume > 0:
			cmd.preMuteVolume = status.Volume
			cmd.volume = 0
		case cmd.sign != 0:
			cmd.volume *= cmd.sign
			cmd.volume = status.Volume + cmd.volume
		}

		return client.SetVolume(cmd.volume)

	default:
		return fmt.Errorf("Unknown input '%s', expected END", t.String())
	}

	return nil
}
