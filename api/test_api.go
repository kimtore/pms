package api

import (
	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/message"
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/style"
)

type testAPI struct {
	messages  chan message.Message
	options   *options.Options
	song      *song.Song
	songlist  songlist.Songlist
	clipboard songlist.Songlist
}

func createTestSong() *song.Song {
	s := song.New()
	s.SetTags(mpd.Attrs{
		"artist": "foo",
		"title":  "bar",
	})
	return s
}

func NewTestAPI() API {
	return &testAPI{
		clipboard: songlist.New(),
		messages:  make(chan message.Message, 1024),
		options:   options.New(),
		song:      createTestSong(),
		songlist:  songlist.New(),
	}
}

func (api *testAPI) Clipboard() songlist.Songlist {
	return api.clipboard
}

func (api *testAPI) Db() *db.Instance {
	return nil // FIXME
}

func (api *testAPI) Exec(cmd string) error {
	panic("not implemented")
}

func (api *testAPI) Multibar() *multibar.Multibar {
	panic("not implemented")
}

func (api *testAPI) Library() *songlist.Library {
	return nil // FIXME
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
	return api.song
}

func (api *testAPI) Songlist() songlist.Songlist {
	return api.songlist
}

func (api *testAPI) Songlists() []songlist.Songlist {
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
