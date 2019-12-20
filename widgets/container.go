package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/multibar"
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
	api      api.API
	events   chan tcell.Event
	multibar *multibar.Multibar
	screen   tcell.Screen
	widgets  widgets
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
		api:    a,
		events: make(chan tcell.Event, 1024),
		screen: screen,
	}, nil
}

func (app *Application) Init() {

	app.widgets.console = views.NewTextArea()

	app.widgets.topbar = NewTopbar()

	app.widgets.multibar = NewMultibarWidget(app.api, app.multibar)

	app.widgets.layout = views.NewBoxLayout(views.Vertical)
	app.widgets.layout.AddWidget(app.widgets.topbar, 1)
	app.widgets.layout.AddWidget(app.widgets.console, 2)
	// app.widgets.layout.AddWidget(app.widgets.songlist, 2)
	app.widgets.layout.AddWidget(app.widgets.multibar, 0)
	app.widgets.layout.SetView(app.screen)
}

func (app *Application) HandleEvent(ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventResize:
		cols, rows := e.Size()
		log.Debugf("terminal resize: %dx%d", cols, rows)
		app.screen.Sync()
		app.widgets.console.Resize()
		return true
	case *tcell.EventKey:
		return false
	default:
		log.Debugf("unrecognized input event: %T %+v", e, e)
		return false
	}
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
