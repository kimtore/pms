package widgets

import (
	console2 "github.com/ambientsound/pms/console"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	log "github.com/sirupsen/logrus"
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
		log.Tracef("keypress: name=%v key=%v modifiers=%v", e.Name(), e.Key(), e.Modifiers())
	case *tcell.EventResize:
		cols, rows := e.Size()
		log.Tracef("terminal resize: %dx%d", cols, rows)
		app.screen.Sync()
		app.widgets.console.Resize()
	default:
		log.Tracef("unrecognized input event: %T %+v", e, e)
	}

	app.widgets.console.HandleEvent(ev)

	return true
}

func (app *Application) Draw() {
	app.widgets.console.SetLines(console2.LogLines)
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
