package commands

import (
	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/message"
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/widgets"
)

// API defines a set of commands that should be available to commands run
// through the command-line interface.
type API interface {
	// Index returns the current Bleve search index, or nil if the search index is not available.
	Index() *index.Index

	// Message sends a message to the user through the statusbar.
	Message(message.Message)

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

	// QuitSignal can be sent to in order to shut down PMS.
	QuitSignal() chan int

	// Song returns the currently playing song, or nil if no song is loaded.
	// Note that the song might be stopped, and the play/pause/stop status should
	// be checked using PlayerStatus().
	Song() *song.Song

	// SonglistWidget returns the songlist widget.
	SonglistWidget() *widgets.SonglistWidget

	// Styles returns the current stylesheet.
	Styles() widgets.StyleMap

	// UI returns the global UI object.
	UI() *widgets.UI
}

type baseAPI struct {
	eventList      chan int
	eventMessage   chan message.Message
	index          func() *index.Index
	mpdClient      func() *mpd.Client
	multibar       *widgets.MultibarWidget
	options        *options.Options
	playerStatus   func() pms_mpd.PlayerStatus
	queue          func() *songlist.Queue
	quitSignal     chan int
	song           func() *song.Song
	songlistWidget func() *widgets.SonglistWidget
	styles         widgets.StyleMap
	ui             *widgets.UI
}

func BaseAPI(
	eventList chan int,
	eventMessage chan message.Message,
	index func() *index.Index,
	mpdClient func() *mpd.Client,
	multibar *widgets.MultibarWidget,
	options *options.Options,
	playerStatus func() pms_mpd.PlayerStatus,
	queue func() *songlist.Queue,
	quitSignal chan int,
	song func() *song.Song,
	songlistWidget func() *widgets.SonglistWidget,
	styles widgets.StyleMap,
	ui *widgets.UI,

) API {
	return &baseAPI{
		eventList:      eventList,
		eventMessage:   eventMessage,
		index:          index,
		mpdClient:      mpdClient,
		multibar:       multibar,
		options:        options,
		playerStatus:   playerStatus,
		queue:          queue,
		quitSignal:     quitSignal,
		song:           song,
		songlistWidget: songlistWidget,
		styles:         styles,
		ui:             ui,
	}
}

func (b *baseAPI) Index() *index.Index {
	return b.index()
}

func (b *baseAPI) Message(msg message.Message) {
	b.eventMessage <- msg
}

func (b *baseAPI) MpdClient() *mpd.Client {
	return b.mpdClient()
}

func (b *baseAPI) Multibar() *widgets.MultibarWidget {
	return b.multibar
}

func (b *baseAPI) Options() *options.Options {
	return b.options
}

func (b *baseAPI) PlayerStatus() pms_mpd.PlayerStatus {
	return b.playerStatus()
}

func (b *baseAPI) Queue() *songlist.Queue {
	return b.queue()
}

func (b *baseAPI) QuitSignal() chan int {
	return b.quitSignal
}

func (b *baseAPI) Song() *song.Song {
	return b.song()
}

func (b *baseAPI) SonglistWidget() *widgets.SonglistWidget {
	return b.songlistWidget()
}

func (b *baseAPI) Styles() widgets.StyleMap {
	return b.styles
}

func (b *baseAPI) UI() *widgets.UI {
	return b.ui
}
