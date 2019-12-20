package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/constants"
	"github.com/ambientsound/pms/log"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type widgets struct {
	layout   *views.BoxLayout
	console  *views.TextArea
	topbar   *Topbar
	multibar *MultibarWidget
	songlist *SonglistWidget
}

type Application struct {
	screen  tcell.Screen
	events  chan tcell.Event
	widgets widgets
	api     api.API
}

var _ tcell.EventHandler = &Application{}

func NewApplication(a api.API) (*Application, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	screen.Clear()
	screen.Show()

	return &Application{
		screen: screen,
		events: make(chan tcell.Event, 1024),
		api:    a,
	}, nil
}

func (app *Application) Init() {

	app.widgets.console = views.NewTextArea()

	app.widgets.topbar = NewTopbar()

	app.widgets.multibar = NewMultibarWidget(app.api, app.Events())

	app.widgets.layout = views.NewBoxLayout(views.Vertical)
	app.widgets.layout.AddWidget(app.widgets.topbar, 1)
	app.widgets.layout.AddWidget(app.widgets.console, 2)
	// app.widgets.layout.AddWidget(app.widgets.songlist, 2)
	app.widgets.layout.AddWidget(app.widgets.multibar, 0)
	app.widgets.layout.SetView(app.screen)
}

func (app *Application) HandleEvent(ev tcell.Event) bool {
	if app.widgets.multibar.HandleEvent(ev) {
		return true
	}

	switch e := ev.(type) {
	case *tcell.EventKey:
		log.Debugf("keypress: name=%v key=%v modifiers=%v", e.Name(), e.Key(), e.Modifiers())
		return false
	case *tcell.EventResize:
		cols, rows := e.Size()
		log.Debugf("terminal resize: %dx%d", cols, rows)
		app.screen.Sync()
		app.widgets.console.Resize()
		return true
	default:
		log.Debugf("unrecognized input event: %T %+v", e, e)
		return true
	}
}

func (app *Application) SetInputMode(mode constants.InputMode) {
	app.widgets.multibar.SetMode(mode)
}

func (app *Application) Draw() {
	app.widgets.console.SetLines(log.Lines())
	app.widgets.layout.Draw()
	app.screen.Show()
}

func (app *Application) Poll() {
	for {
		app.events <- app.screen.PollEvent()
	}
}

func (app *Application) Events() <-chan tcell.Event {
	return app.events
}

func (app *Application) Finish() {
	app.screen.Fini()
}
