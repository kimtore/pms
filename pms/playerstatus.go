package pms

import "time"

// PlayerStatus contains information about MPD's player status.
type PlayerStatus struct {
	Audio          string
	Bitrate        int
	Consume        bool
	Elapsed        float64
	Err            string
	MixRampDB      float64
	Playlist       int
	PlaylistLength int
	Random         bool
	Repeat         bool
	Single         bool
	Song           int
	SongID         int
	State          string
	Time           int
	Volume         int

	updateTime time.Time
}

// Strings found in the PlayerStatus.State variable.
const (
	StatePlay    string = "play"
	StateStop    string = "stop"
	StatePause   string = "pause"
	StateUnknown string = "unknown"
)

func (p *PlayerStatus) SetTime() {
	p.updateTime = time.Now()
}

func (p *PlayerStatus) Since() time.Duration {
	return time.Since(p.updateTime)
}

func (p *PlayerStatus) Tick() {
	if p.State != StatePlay {
		return
	}
	diff := p.Since()
	p.SetTime()
	p.Elapsed += diff.Seconds()
}
