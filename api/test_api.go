package api

import (
	"fmt"
	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/player"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/spotify/library"
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/ambientsound/pms/style"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify"
)

type testAPI struct {
	messages  chan message.Message
	list      list.List
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
		list:      list.New(),
		messages:  make(chan message.Message, 1024),
		song:      createTestSong(),
		songlist:  songlist.New(),
	}
}

func (api *testAPI) Authenticate() error {
	return nil
}

func (api *testAPI) Clipboard() songlist.Songlist {
	return api.clipboard
}

func (api *testAPI) Db() *db.List {
	return nil // FIXME
}

func (api *testAPI) Exec(cmd string) error {
	panic("not implemented")
}

func (api *testAPI) Multibar() *multibar.Multibar {
	panic("not implemented")
}

func (api *testAPI) List() list.List {
	return api.list
}

func (api *testAPI) Library() *spotify_library.List {
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

func (api *testAPI) Options() Options {
	return viper.GetViper()
}

func (api *testAPI) PlayerStatus() player.State {
	return player.State{}
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

func (api *testAPI) SetList(lst list.List) {
}

func (api *testAPI) Spotify() (*spotify.Client, error) {
	return nil, fmt.Errorf("no spotify")
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

func (api *testAPI) Styles() style.Stylesheet {
	return nil // FIXME
}

func (api *testAPI) Tracklist() *spotify_tracklist.List {
	return nil // FIXME
}

func (api *testAPI) UI() UI {
	return nil // FIXME
}
