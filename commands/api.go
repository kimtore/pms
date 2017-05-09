package commands

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
	"github.com/ambientsound/pms/widgets"
)

// API defines a set of commands that should be available to commands run
// through the command-line interface.
type API interface {
	// Index returns the current Bleve search index, or nil if the search index is not available.
	Index() *index.Index

	// ListChanged notifies the UI that the current songlist has changed.
	ListChanged()

	// OptionChanged notifies that an option has been changed.
	OptionChanged(string)

	// Message sends a message to the user through the statusbar.
	Message(string, ...interface{})

	// MpdClient returns the current MPD client, which is confirmed to be alive. If the MPD connection is not working, nil is returned.
	MpdClient() *mpd.Client

	// Multibar returns the Multibar widget.
	Multibar() *widgets.MultibarWidget

	// Options returns PMS' global options.
	Options() *options.Options

	// PlayerStatus returns the current MPD player status.
	PlayerStatus() pms_mpd.PlayerStatus

	// Queue returns MPD's song queue.
	Queue() *songlist.Queue

	// Quit shuts down PMS.
	Quit()

	// Sequencer returns a pointer to the key sequencer that receives key events.
	Sequencer() *keys.Sequencer

	// Song returns the currently playing song, or nil if no song is loaded.
	// Note that the song might be stopped, and the play/pause/stop status should
	// be checked using PlayerStatus().
	Song() *song.Song

	// SonglistWidget returns the songlist widget.
	SonglistWidget() *widgets.SonglistWidget

	// Styles returns the current stylesheet.
	Styles() style.Stylesheet

	// UI returns the global UI object.
	UI() *widgets.UI
}

type baseAPI struct {
	eventList      chan int
	eventMessage   chan message.Message
	eventOption    chan string
	index          func() *index.Index
	mpdClient      func() *mpd.Client
	multibar       *widgets.MultibarWidget
	options        *options.Options
	playerStatus   func() pms_mpd.PlayerStatus
	queue          func() *songlist.Queue
	quitSignal     chan int
	sequencer      *keys.Sequencer
	song           func() *song.Song
	songlistWidget func() *widgets.SonglistWidget
	styles         style.Stylesheet
	ui             *widgets.UI
}

func BaseAPI(
	eventList chan int,
	eventMessage chan message.Message,
	eventOption chan string,
	index func() *index.Index,
	mpdClient func() *mpd.Client,
	multibar *widgets.MultibarWidget,
	options *options.Options,
	playerStatus func() pms_mpd.PlayerStatus,
	queue func() *songlist.Queue,
	quitSignal chan int,
	sequencer *keys.Sequencer,
	song func() *song.Song,
	songlistWidget func() *widgets.SonglistWidget,
	styles style.Stylesheet,
	ui *widgets.UI,

) API {
	return &baseAPI{
		eventList:      eventList,
		eventMessage:   eventMessage,
		eventOption:    eventOption,
		index:          index,
		mpdClient:      mpdClient,
		multibar:       multibar,
		options:        options,
		playerStatus:   playerStatus,
		queue:          queue,
		quitSignal:     quitSignal,
		sequencer:      sequencer,
		song:           song,
		songlistWidget: songlistWidget,
		styles:         styles,
		ui:             ui,
	}
}

func (api *baseAPI) Index() *index.Index {
	return api.index()
}

func (api *baseAPI) ListChanged() {
	api.eventList <- 0
}

func (api *baseAPI) Message(fmt string, a ...interface{}) {
	api.eventMessage <- message.Format(fmt, a...)
}

func (api *baseAPI) MpdClient() *mpd.Client {
	return api.mpdClient()
}

func (api *baseAPI) Multibar() *widgets.MultibarWidget {
	return api.multibar
}

func (api *baseAPI) OptionChanged(key string) {
	api.eventOption <- key
}

func (api *baseAPI) Options() *options.Options {
	return api.options
}

func (api *baseAPI) PlayerStatus() pms_mpd.PlayerStatus {
	return api.playerStatus()
}

func (api *baseAPI) Queue() *songlist.Queue {
	return api.queue()
}

func (api *baseAPI) Quit() {
	api.quitSignal <- 0
}

func (api *baseAPI) Sequencer() *keys.Sequencer {
	return api.sequencer
}

func (api *baseAPI) Song() *song.Song {
	return api.song()
}

func (api *baseAPI) SonglistWidget() *widgets.SonglistWidget {
	return api.songlistWidget()
}

func (api *baseAPI) Styles() style.Stylesheet {
	return api.styles
}

func (api *baseAPI) UI() *widgets.UI {
	return api.ui
}
