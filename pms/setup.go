package pms

import (
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/topbar"
	"github.com/ambientsound/pms/widgets"
)

func New() *PMS {
	pms := &PMS{}

	pms.EventIndex = make(chan int)
	pms.EventLibrary = make(chan int)
	pms.EventList = make(chan int, 1024)
	pms.EventMessage = make(chan message.Message, 1024)
	pms.EventPlayer = make(chan int)
	pms.EventOption = make(chan string, 1024)
	pms.EventQueue = make(chan int)
	pms.QuitSignal = make(chan int, 1)
	pms.stylesheet = make(style.Stylesheet)

	pms.Queue = songlist.NewQueue(pms.CurrentMpdClient)
	pms.Library = songlist.NewLibrary()

	pms.Options = options.New()
	pms.Options.AddDefaultOptions()

	pms.Sequencer = keys.NewSequencer()

	pms.clipboards = make(map[string]songlist.Songlist)
	pms.clipboards["default"] = songlist.New()

	pms.setupUI()

	pms.CLI = input.NewCLI(pms.API())

	return pms
}

// setupAPI creates an API object
func (pms *PMS) API() api.API {
	return api.BaseAPI(
		pms.Clipboard,
		pms.EventList,
		pms.EventMessage,
		pms.EventOption,
		pms.CurrentIndex,
		pms.CurrentMpdClient,
		pms.Multibar,
		pms.Options,
		pms.CurrentPlayerStatus,
		pms.CurrentQueue,
		pms.QuitSignal,
		pms.Sequencer,
		pms.CurrentSong,
		pms.CurrentSonglistWidget,
		pms.Stylesheet(),
		pms.UI,
	)
}

func (pms *PMS) setupUI() {
	timer := time.Now()
	pms.ui = widgets.NewUI(pms.API())
	pms.ui.Start()
	pms.ui.Songlist.AddSonglist(pms.Queue)
	pms.ui.Songlist.AddSonglist(pms.Library)
	pms.ui.Songlist.SetSonglist(pms.Queue)

	console.Log("UI initialized in %s", time.Since(timer).String())
}

func (pms *PMS) setupTopbar() {
	config := pms.Options.StringValue("topbar")
	matrix, err := topbar.Parse(pms.API(), config)
	if err == nil {
		pms.ui.Topbar.SetMatrix(matrix)
	} else {
		pms.Error("Error in topbar configuration: %s", err)
	}
}
