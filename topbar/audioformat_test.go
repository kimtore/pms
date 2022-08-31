package topbar_test

import (
	"fmt"
	"testing"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/topbar"
)

func TestAudioformats(t *testing.T) {
	var testcases = []struct {
		in, param, want string
	}{
		{"(44100:16:2)", "", "(44100:16:2)"},
		{"(192000:24:2)", "", "(192000:24:2)"},
		{"(192000:24:2)", "channels", "2 ch"},
		{"(192000:24:2)", "resolution", "24 bit"},
		{"(192000:24:2)", "samplerate", "192 kHz"},
		{"(dsd128:5)", "channels", "5 ch"},
		{"(dsd128:5)", "resolution", "1 bit"},
		{"(dsd128:5)", "samplerate", fmt.Sprintf("%v MHz", 128*44100.0/1e6)},
	}
	for _, tc := range testcases {
		api := api.NewTestAPI()
		p := api.PlayerStatus()
		p.Audio = tc.in
		api.Db().SetPlayerStatus(p)
		af := topbar.NewAudioformat(api, tc.param)
		gotstring, _ := af.Text()
		if gotstring != tc.want {
			t.Errorf("Wrong audioformat string: got %v, want %v", gotstring, tc.want)
		}
	}
}
