package player

import (
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/zmb3/spotify"
	"time"
)

// State contains information about MPD's player status.
type State struct {
	spotify.PlayerState

	CreateTime         time.Time
	ProgressPercentage float64
	TrackRow           list.Row
	updateTime         time.Time
}

func NewState(state spotify.PlayerState) State {
	row := list.Row{}
	if state.Item != nil {
		row = spotify_tracklist.Row(*state.Item)
	}
	return State{
		PlayerState: state,
		CreateTime:  time.Now(),
		TrackRow:    row,
		updateTime:  time.Now(),
	}
}

const (
	StatePlay    string = "play"
	StateStop    string = "stop"
	StatePause   string = "pause"
	StateUnknown string = "unknown"
)

func (p *State) SetTime() {
	p.updateTime = time.Now()
}

func (p *State) Since() time.Duration {
	return time.Since(p.updateTime)
}

func (p State) State() string {
	// FIXME
	if p.Playing {
		return StatePlay
	}
	if p.Item == nil {
		return StateStop
	}
	return StatePause
}

func (p State) percentage() float64 {
	if p.Item == nil {
		return p.ProgressPercentage
	} else if p.Progress == 0 {
		return 0.0
	} else {
		return float64(p.Progress) / float64(p.Item.Duration)
	}
}

func (p State) Tick() State {
	if !p.Playing {
		return p
	}
	diff := p.Since()
	p.SetTime()
	p.Progress += int(diff.Milliseconds())
	p.ProgressPercentage = p.percentage()

	return p
}
