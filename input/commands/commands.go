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

type Base struct {
	CurrentIndex        func() *index.Index
	CurrentPlayerStatus func() pms_mpd.PlayerStatus
	CurrentQueue        func() *songlist.Queue
	CurrentSong         func() *song.Song
	EventList           chan int
	EventMessage        chan message.Message
	MpdClient           func() *mpd.Client
	Multibar            *widgets.MultibarWidget
	Options             *options.Options
	QuitSignal          chan int
	SonglistWidget      func() *widgets.SonglistWidget
	Styles              widgets.StyleMap
	Ui                  *widgets.UI
}
