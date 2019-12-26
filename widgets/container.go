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
	active   views.Widget
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
	app.Widgets.Topbar = NewTopbar()
	app.Widgets.table = NewTable(app.api)
	app.Widgets.multibar = NewMultibarWidget(app.api)

	app.Widgets.layout = views.NewBoxLayout(views.Vertical)
	app.Widgets.layout.AddWidget(app.Widgets.Topbar, 0)
	app.Widgets.layout.AddWidget(app.Widgets.table, 1)
	app.Widgets.layout.AddWidget(app.Widgets.multibar, 0)
	app.Widgets.layout.SetView(app.screen)

	app.Widgets.active = app.Widgets.table

	app.ActivateWindow(api.WindowLogs)
}

func (app *Application) HandleEvent(ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventResize:
		cols, rows := e.Size()
		log.Debugf("terminal resize: %dx%d", cols, rows)
		app.screen.Sync()
		app.Widgets.layout.Resize()
		app.Widgets.layout.SetView(app.screen)
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

func (app *Application) ActivateWindow(window api.Window) {
	var widget views.Widget

	switch window {
	case api.WindowLogs:
		widget = app.Widgets.table
		app.Widgets.table.SetList(log.List(log.InfoLevel))
		app.Widgets.table.SetColumns([]string{"timestamp", "logLevel", "logMessage"})
	case api.WindowMusic:
		widget = app.Widgets.table
		app.Widgets.table.SetList(nil)
	case api.WindowPlaylists:
		widget = app.Widgets.table
		app.Widgets.table.SetList(nil)
	default:
		panic("widget not implemented")
	}

	log.Debugf("want to activate widget %#v", widget)
	log.Debugf("first deactivating widget %#v", app.Widgets.active)

	app.Widgets.layout.RemoveWidget(app.Widgets.active)
	app.Widgets.layout.InsertWidget(1, widget, 1.0)

	app.Widgets.active = widget
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
