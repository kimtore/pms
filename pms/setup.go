package pms

import (
	"strings"
	"time"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/commands"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/widgets"
)

func New() *PMS {
	pms := &PMS{}

	pms.EventError = make(chan string, 16)
	pms.EventIndex = make(chan int)
	pms.EventLibrary = make(chan int)
	pms.EventQueue = make(chan int)
	pms.EventMessage = make(chan string, 16)
	pms.EventPlayer = make(chan int)
	pms.QuitSignal = make(chan int, 1)

	pms.Options = options.New()
	pms.Options.AddDefaultOptions()

	pms.Sequencer = keys.NewSequencer()

	pms.setupUI()
	pms.setupCLI()
	pms.readDefaultConfiguration()

	return pms
}

// SetupCLI instantiates the different commands PMS understands, such as set; bind; etc.
func (pms *PMS) setupCLI() {
	pms.CLI = input.NewCLI()
	pms.CLI.Register("bind", commands.NewBind(pms.Sequencer))
	pms.CLI.Register("cursor", commands.NewCursor(pms.UI.Songlist))
	pms.CLI.Register("inputmode", commands.NewInputMode(pms.UI.Multibar))
	pms.CLI.Register("play", commands.NewPlay(pms.UI.Songlist, pms.CurrentMpdClient))
	pms.CLI.Register("q", commands.NewQuit(pms.QuitSignal))
	pms.CLI.Register("quit", commands.NewQuit(pms.QuitSignal))
	pms.CLI.Register("redraw", commands.NewRedraw(pms.UI.App))
	pms.CLI.Register("se", commands.NewSet(pms.Options))
	pms.CLI.Register("set", commands.NewSet(pms.Options))
}

func (pms *PMS) setupUI() {
	timer := time.Now()
	pms.UI = widgets.NewUI(pms.Options)
	pms.UI.Start()
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
