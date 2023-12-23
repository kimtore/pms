package api

import (
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/message"
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/style"
	"github.com/fhs/gompd/v2/mpd"
)

type testAPI struct {
	messages  chan message.Message
	options   *options.Options
	song      *song.Song
	songlist  songlist.Songlist
	clipboard songlist.Songlist
	db        *db.Instance
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
		db:        db.New(),
	}
}

func (api *testAPI) Clipboard() songlist.Songlist {
	return api.clipboard
}

func (api *testAPI) Db() *db.Instance {
	return api.db
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

func (api *testAPI) Multibar() MultibarWidget {
	return nil // FIXME
}

func (api *testAPI) OptionChanged(key string) {
	// FIXME
}

func (api *testAPI) Options() *options.Options {
	return api.options
}

func (api *testAPI) PlayerStatus() pms_mpd.PlayerStatus {
	return api.db.PlayerStatus()
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

// SetPlayerStatus sets the player status struct to the provided input
// value, allowing to construct test cases that depend on a particular
// player status
//func (api *testAPI) SetPlayerStatus(p pms_mpd.PlayerStatus) {
//api.status = p
//}

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
