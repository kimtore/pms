package pms

import (
	"strings"
	"time"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/commands"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/widgets"
)

func New() *PMS {
	pms := &PMS{}

	pms.EventIndex = make(chan int)
	pms.EventLibrary = make(chan int)
	pms.EventList = make(chan int, 1024)
	pms.EventMessage = make(chan message.Message, 1024)
	pms.EventPlayer = make(chan int)
	pms.EventQueue = make(chan int)
	pms.QuitSignal = make(chan int, 1)

	pms.Queue = songlist.NewQueue(pms.CurrentMpdClient)
	pms.Library = songlist.NewLibrary()

	pms.Options = options.New()
	pms.Options.AddDefaultOptions()

	pms.Sequencer = keys.NewSequencer()

	pms.setupUI()
	pms.setupCLI()
	pms.readDefaultConfiguration()

	return pms
}

// setupCLI instantiates the different commands PMS understands, such as set; bind; etc.
func (pms *PMS) setupCLI() {
	resources := commands.Base{
		CurrentIndex:        pms.CurrentIndex,
		CurrentPlayerStatus: pms.CurrentPlayerStatus,
		CurrentQueue:        pms.CurrentQueue,
		CurrentSong:         pms.CurrentSong,
		SonglistWidget:      pms.UI.CurrentSonglistWidget,
		EventList:           pms.EventList,
		EventMessage:        pms.EventMessage,
		MpdClient:           pms.CurrentMpdClient,
		Multibar:            pms.UI.Multibar,
		Options:             pms.Options,
		QuitSignal:          pms.QuitSignal,
		Styles:              pms.UI.Styles,
		Ui:                  pms.UI,
	}

	pms.CLI = input.NewCLI(resources)
	pms.CLI.Registerf("add", commands.NewAdd)
	pms.CLI.Register("bind", commands.NewBind(pms.Sequencer))
	pms.CLI.Register("cursor", commands.NewCursor(pms.UI.CurrentSonglistWidget, pms.CurrentSong))
	pms.CLI.Register("inputmode", commands.NewInputMode(pms.UI.Multibar))
	pms.CLI.Register("isolate", commands.NewIsolate(pms.EventMessage, pms.UI.CurrentSonglistWidget, pms.CurrentIndex, pms.Options))
	pms.CLI.Register("list", commands.NewList(pms.UI))
	pms.CLI.Register("next", commands.NewNext(pms.CurrentMpdClient))
	pms.CLI.Register("pause", commands.NewPause(pms.CurrentMpdClient, pms.CurrentPlayerStatus))
	pms.CLI.Register("play", commands.NewPlay(pms.UI.CurrentSonglistWidget, pms.CurrentMpdClient))
	pms.CLI.Register("prev", commands.NewPrevious(pms.CurrentMpdClient))
	pms.CLI.Register("previous", commands.NewPrevious(pms.CurrentMpdClient))
	pms.CLI.Register("q", commands.NewQuit(pms.QuitSignal))
	pms.CLI.Register("quit", commands.NewQuit(pms.QuitSignal))
	pms.CLI.Register("redraw", commands.NewRedraw(pms.UI.App))
	pms.CLI.Register("remove", commands.NewRemove(pms.UI.CurrentSonglistWidget, pms.EventList))
	pms.CLI.Register("se", commands.NewSet(pms.Options, pms.EventMessage))
	pms.CLI.Register("select", commands.NewSelect(pms.UI.CurrentSonglistWidget))
	pms.CLI.Register("set", commands.NewSet(pms.Options, pms.EventMessage))
	pms.CLI.Register("sort", commands.NewSort(pms.UI.CurrentSonglistWidget, pms.Options))
	pms.CLI.Register("stop", commands.NewStop(pms.CurrentMpdClient))
	pms.CLI.Register("style", commands.NewStyle(pms.UI.Styles))
	pms.CLI.Register("volume", commands.NewVolume(pms.CurrentMpdClient, pms.CurrentPlayerStatus))
}

func (pms *PMS) setupUI() {
	timer := time.Now()
	pms.UI = widgets.NewUI(pms.Options)
	pms.UI.Start()
	pms.UI.Songlist.AddSonglist(pms.Queue)
	pms.UI.Songlist.AddSonglist(pms.Library)
	pms.UI.Songlist.SetSonglist(pms.Queue)
	console.Log("UI initialized in %s", time.Since(timer).String())
}

func (pms *PMS) readDefaultConfiguration() {
	lines := strings.Split(options.Defaults, "\n")
	for _, line := range lines {
		err := pms.CLI.Execute(line)
		if err != nil {
			console.Log("Error while reading default configuration: %s", err)
		}
	}
}
