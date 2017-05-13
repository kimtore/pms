package pms

import (
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
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

	pms.setupUI()
	pms.setupCLI()

	return pms
}

// setupAPI creates an API object
func (pms *PMS) API() api.API {
	return api.BaseAPI(
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

// setupCLI instantiates the different commands PMS understands, such as set; bind; etc.
func (pms *PMS) setupCLI() {
	pms.CLI = input.NewCLI(pms.API())
	pms.CLI.Register("add", commands.NewAdd)
	pms.CLI.Register("bind", commands.NewBind)
	pms.CLI.Register("cursor", commands.NewCursor)
	pms.CLI.Register("inputmode", commands.NewInputMode)
	pms.CLI.Register("isolate", commands.NewIsolate)
	pms.CLI.Register("list", commands.NewList)
	pms.CLI.Register("next", commands.NewNext)
	pms.CLI.Register("pause", commands.NewPause)
	pms.CLI.Register("play", commands.NewPlay)
	pms.CLI.Register("prev", commands.NewPrevious)
	pms.CLI.Register("previous", commands.NewPrevious)
	pms.CLI.Register("print", commands.NewPrint)
	pms.CLI.Register("q", commands.NewQuit)
	pms.CLI.Register("quit", commands.NewQuit)
	pms.CLI.Register("redraw", commands.NewRedraw)
	pms.CLI.Register("remove", commands.NewRemove)
	pms.CLI.Register("se", commands.NewSet)
	pms.CLI.Register("select", commands.NewSelect)
	pms.CLI.Register("set", commands.NewSet)
	pms.CLI.Register("sort", commands.NewSort)
	pms.CLI.Register("stop", commands.NewStop)
	pms.CLI.Register("style", commands.NewStyle)
	pms.CLI.Register("volume", commands.NewVolume)
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
