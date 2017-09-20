package pms

import (
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/db"
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

	pms.database = db.New()

	pms.EventLibrary = make(chan int, 1024)
	pms.EventList = make(chan int, 1024)
	pms.EventMessage = make(chan message.Message, 1024)
	pms.EventPlayer = make(chan int, 1024)
	pms.EventOption = make(chan string, 1024)
	pms.EventQueue = make(chan int, 1024)
	pms.QuitSignal = make(chan int, 1)
	pms.stylesheet = make(style.Stylesheet)

	pms.database.SetQueue(songlist.NewQueue(pms.CurrentMpdClient))
	pms.database.SetLibrary(songlist.NewLibrary())

	pms.Options = options.New()
	pms.Options.AddDefaultOptions()

	pms.Sequencer = keys.NewSequencer()

	pms.setupUI()

	pms.CLI = input.NewCLI(pms.API())

	return pms
}

// setupAPI creates an API object
func (pms *PMS) API() api.API {
	return api.BaseAPI(
		pms.Database,
		pms.EventList,
		pms.EventMessage,
		pms.EventOption,
		pms.database.Library,
		pms.CurrentMpdClient,
		pms.Multibar,
		pms.Options,
		pms.CurrentPlayerStatus,
		pms.database.Queue,
		pms.QuitSignal,
		pms.Sequencer,
		pms.database.CurrentSong,
		pms.CurrentSonglistWidget,
		pms.Stylesheet(),
		pms.UI,
	)
}

func (pms *PMS) setupUI() {
	timer := time.Now()
	queue := pms.database.Queue()
	pms.ui = widgets.NewUI(pms.API())
	pms.ui.Start()
	pms.database.Panel().Add(queue)
	pms.database.Panel().Add(pms.database.Library())
	pms.database.Panel().Activate(queue)

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
