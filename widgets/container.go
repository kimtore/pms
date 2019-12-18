package widgets

import (
	"github.com/ambientsound/pms/log"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type widgets struct {
	console *views.TextArea
}

type Application struct {
	screen  tcell.Screen
	events  chan tcell.Event
	widgets widgets
}

var _ tcell.EventHandler = &Application{}

func NewApplication() (*Application, error) {
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

	console := views.NewTextArea()
	console.SetView(screen)

	return &Application{
		screen: screen,
		events: make(chan tcell.Event, 1024),
		widgets: widgets{
			console: console,
		},
	}, nil
}

func (app *Application) HandleEvent(ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventKey:
		log.Debugf("keypress: name=%v key=%v modifiers=%v", e.Name(), e.Key(), e.Modifiers())
	case *tcell.EventResize:
		cols, rows := e.Size()
		log.Debugf("terminal resize: %dx%d", cols, rows)
		app.screen.Sync()
		app.widgets.console.Resize()
	default:
		log.Debugf("unrecognized input event: %T %+v", e, e)
	}

	app.widgets.console.HandleEvent(ev)

	return true
}

func (app *Application) Draw() {
	app.widgets.console.SetLines(log.Lines())
	app.widgets.console.Draw()
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
