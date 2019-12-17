package widgets

import (
	"github.com/gdamore/tcell"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	screen tcell.Screen
	events chan tcell.Event
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

	return &Application{
		screen: screen,
		events: make(chan tcell.Event, 1024),
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
	default:
		log.Tracef("unrecognized input event: %T %+v", e, e)
	}
	return true
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
