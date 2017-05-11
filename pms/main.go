package pms

import (
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/message"
)

func (pms *PMS) Main() {
	for {
		select {
		case <-pms.QuitSignal:
			pms.handleQuitSignal()
			return
		case <-pms.EventLibrary:
			pms.handleEventLibrary()
		case <-pms.EventQueue:
			pms.handleEventQueue()
		case <-pms.EventIndex:
			pms.handleEventIndex()
		case <-pms.EventList:
			pms.handleEventList()
		case <-pms.EventPlayer:
			pms.handleEventPlayer()
		case key := <-pms.EventOption:
			pms.handleEventOption(key)
		case msg := <-pms.EventMessage:
			pms.handleEventMessage(msg)
		case ev := <-pms.ui.EventKeyInput:
			pms.KeyInput(ev)
		case s := <-pms.ui.EventInputCommand:
			pms.Execute(s)
		}

		// Draw missing parts after every iteration
		pms.ui.App.PostFunc(func() {
			pms.ui.App.Update()
		})
	}
}

func (pms *PMS) handleQuitSignal() {
	console.Log("Received quit signal, exiting.")
	pms.ui.Quit()
}

func (pms *PMS) handleEventLibrary() {
	console.Log("Song library updated in MPD, assigning to UI")
	pms.ui.App.PostFunc(func() {
		pms.ui.Songlist.ReplaceSonglist(pms.Library)
	})
}

func (pms *PMS) handleEventQueue() {
	console.Log("Queue updated in MPD, assigning to UI")
	pms.ui.App.PostFunc(func() {
		pms.ui.Songlist.ReplaceSonglist(pms.Queue)
	})
}

func (pms *PMS) handleEventIndex() {
	console.Log("Search index updated, assigning to UI")
	pms.ui.App.PostFunc(func() {
		pms.ui.SetIndex(pms.Index)
	})
}

func (pms *PMS) handleEventList() {
	console.Log("Songlist changed, notifying UI")
	pms.ui.App.PostFunc(func() {
		pms.ui.Songlist.ListChanged()
	})
}

func (pms *PMS) handleEventOption(key string) {
	console.Log("Option '%s' has been changed", key)
	switch key {
	case "topbar":
		pms.setupTopbar()
	case "columns":
		pms.ui.App.PostFunc(func() {
			pms.ui.Songlist.ListChanged()
		})
	}
}

func (pms *PMS) handleEventPlayer() {
}

func (pms *PMS) handleEventMessage(msg message.Message) {
	message.Log(msg)
	pms.ui.App.PostFunc(func() {
		pms.ui.Multibar.SetMessage(msg)
	})
}
