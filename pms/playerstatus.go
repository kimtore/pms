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
