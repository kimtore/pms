package api

import (
	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/message"
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/style"
)

type testAPI struct {
	messages chan message.Message
	options  *options.Options
}

func NewTestAPI() API {
	return &testAPI{
		messages: make(chan message.Message, 1024),
		options:  options.New(),
	}
}

func (api *testAPI) Index() *index.Index {
	return nil
}

func (api *testAPI) ListChanged() {
	// FIXME
}

func (api *testAPI) Message(fmt string, a ...interface{}) {
	api.messages <- message.Format(fmt, a...)
}

func (api *testAPI) MpdClient() *mpd.Client {
	return nil // FIXME
}

func (api *testAPI) Multibar() MultibarWidget {
	return nil // FIXME
}

func (api *testAPI) OptionChanged(key string) {
	// FIXME
}

func (api *testAPI) Options() *options.Options {
	return api.options
}

func (api *testAPI) PlayerStatus() (p pms_mpd.PlayerStatus) {
	return // FIXME
}

func (api *testAPI) Queue() *songlist.Queue {
	return nil // FIXME
}

func (api *testAPI) Quit() {
	return // FIXME
}

func (api *testAPI) Sequencer() *keys.Sequencer {
	return nil // FIXME
}

func (api *testAPI) Song() *song.Song {
	return nil // FIXME
}

func (api *testAPI) SonglistWidget() SonglistWidget {
	return nil // FIXME
}

func (api *testAPI) Styles() style.Stylesheet {
	return nil // FIXME
}

func (api *testAPI) UI() UI {
	return nil // FIXME
}
