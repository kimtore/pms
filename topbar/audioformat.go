package topbar

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
)

// Audioformat draws the current audio format.
type Audioformat struct {
	api    api.API
	aspect string
}

// NewAudioformat returns Audioformat.
func NewAudioformat(a api.API, param string) Fragment {
	return &Audioformat{a, param}
}

// Text implements Fragment.
func (w *Audioformat) Text() (string, string) {
	playerStatus := w.api.PlayerStatus()

	if w.aspect == "" {
		text := fmt.Sprintf("%v", playerStatus.Audio)
		return text, `audioformat`
	}

	audioformat := strings.Split(strings.Trim(playerStatus.Audio, "()"), ":")
	formatmap := make(map[string]string)
	// The audioformat string is defined at
	// https://github.com/MusicPlayerDaemon/MPD/blob/master/src/pcm/AudioFormat.cxx
	switch len(audioformat) {
	case 3:
		// mpd sends a tuple like (44100:16:2)
		kHz, _ := strconv.Atoi(audioformat[0])
		formatmap["samplerate"] = fmt.Sprintf("%v kHz", float64(kHz)/1000.0)
		formatmap["channels"] = fmt.Sprintf("%v ch", audioformat[2])
		switch audioformat[1] {
		case "f":
			formatmap["resolution"] = "float"
		default:
			formatmap["resolution"] = fmt.Sprintf("%v bit", audioformat[1])
		}
		return formatmap[w.aspect], w.aspect
	case 2:
		// mpd sends "dsd<number>" strings
		samplingfactor, _ := strconv.Atoi(audioformat[0][3:])
		formatmap["samplerate"] = fmt.Sprintf("%v MHz", float64(samplingfactor)*44100.0/1e6)
		formatmap["resolution"] = "1 bit"
		formatmap["channels"] = fmt.Sprintf("%v ch", audioformat[1])
		return formatmap[w.aspect], w.aspect
	default:
		// If we end up here, something in mpd has changed
		console.Log("Unsupported audio format string: %s", playerStatus.Audio)
		text := fmt.Sprintf("%v", playerStatus.Audio)
		return text, `audioformat`
	}
}
