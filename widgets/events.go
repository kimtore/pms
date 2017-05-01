package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type eventfulWidget interface {
	views.Widget
	PostEvent(wev views.EventWidget)
}

type widgetEvent struct {
	widget views.Widget
	tcell.EventTime
}

type EventInputChanged struct {
	widgetEvent
}

func PostEventInputChanged(w eventfulWidget) {
	ev := &EventInputChanged{}
	ev.SetWidget(w)
	w.PostEvent(ev)
}

type EventInputFinished struct {
	widgetEvent
}

func PostEventInputFinished(w eventfulWidget) {
	ev := &EventInputFinished{}
	ev.SetWidget(w)
	w.PostEvent(ev)
}

type EventListChanged struct {
	widgetEvent
}

func PostEventListChanged(w eventfulWidget) {
	ev := &EventListChanged{}
	ev.SetWidget(w)
	w.PostEvent(ev)
}

type EventScroll struct {
	widgetEvent
}

func PostEventScroll(w eventfulWidget) {
	ev := &EventScroll{}
	ev.SetWidget(w)
	w.PostEvent(ev)
}

type EventModeSync struct {
	widgetEvent
	InputMode int
}

func PostEventModeSync(w eventfulWidget, inputMode int) {
	ev := &EventModeSync{InputMode: inputMode}
	ev.SetWidget(w)
	w.PostEvent(ev)
}

func (wev *widgetEvent) Widget() views.Widget {
	return wev.widget
}

func (wev *widgetEvent) SetWidget(widget views.Widget) {
	wev.widget = widget
}
