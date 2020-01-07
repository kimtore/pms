package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/style"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type widgets struct {
	layout   *views.BoxLayout
	Topbar   *Topbar
	multibar *Multibar
	table    *Table
}

type Application struct {
	api     api.API
	events  chan tcell.Event
	screen  tcell.Screen
	Widgets widgets
	style.Styled
}

var _ tcell.EventHandler = &Application{}

var _ api.UI = &Application{}

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
	app.Widgets.Topbar = NewTopbar(app.api)
	app.Widgets.table = NewTable(app.api)
	app.Widgets.multibar = NewMultibarWidget(app.api)
	app.Resize()
}

func (app *Application) Resize() {
	app.Widgets.layout = views.NewBoxLayout(views.Vertical)
	app.Widgets.layout.AddWidget(app.Widgets.Topbar, 0)
	app.Widgets.layout.AddWidget(app.Widgets.table, 1)
	app.Widgets.layout.AddWidget(app.Widgets.multibar, 0)
	app.Widgets.layout.SetView(app.screen)
}

func (app *Application) HandleEvent(ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventResize:
		cols, rows := e.Size()
		log.Debugf("terminal resize: %dx%d", cols, rows)
		app.screen.Sync()
		app.Resize()
		return true
	case *tcell.EventKey:
		return false
	default:
		log.Debugf("unrecognized input event: %T %+v", e, e)
		app.Widgets.layout.HandleEvent(ev)
		return false
	}
}

func (app *Application) Draw() {
	app.Widgets.layout.Draw()
	app.updateCursor()
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

func (app *Application) Refresh() {
	app.screen.Sync()
}

func (app *Application) TableWidget() api.TableWidget {
	return app.Widgets.table
}

func (app *Application) updateCursor() {
	switch app.api.Multibar().Mode() {
	case multibar.ModeInput, multibar.ModeSearch:
		_, ymax := app.screen.Size()
		x := app.api.Multibar().Cursor() + 1
		app.screen.ShowCursor(x, ymax-1)
	default:
		app.screen.HideCursor()
	}
}
