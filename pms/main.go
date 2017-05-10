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
		case ev := <-pms.UI.EventKeyInput:
			pms.KeyInput(ev)
		case s := <-pms.UI.EventInputCommand:
			pms.Execute(s)
		}

		// Draw missing parts after every iteration
		pms.UI.App.PostFunc(func() {
			pms.UI.App.Update()
		})
	}
}

func (pms *PMS) handleQuitSignal() {
	console.Log("Received quit signal, exiting.")
	pms.UI.Quit()
}

func (pms *PMS) handleEventLibrary() {
	console.Log("Song library updated in MPD, assigning to UI")
	pms.UI.App.PostFunc(func() {
		pms.UI.Songlist.ReplaceSonglist(pms.Library)
	})
}

func (pms *PMS) handleEventQueue() {
	console.Log("Queue updated in MPD, assigning to UI")
	pms.UI.App.PostFunc(func() {
		pms.UI.Songlist.ReplaceSonglist(pms.Queue)
	})
}

func (pms *PMS) handleEventIndex() {
	console.Log("Search index updated, assigning to UI")
	pms.UI.App.PostFunc(func() {
		pms.UI.SetIndex(pms.Index)
	})
}

func (pms *PMS) handleEventList() {
	console.Log("Songlist changed, notifying UI")
	pms.UI.App.PostFunc(func() {
		pms.UI.Songlist.ListChanged()
	})
}

func (pms *PMS) handleEventOption(key string) {
	console.Log("Option '%s' has been changed", key)
	switch key {
	case "topbar":
		pms.setupTopbar()
	case "columns":
		pms.UI.App.PostFunc(func() {
			pms.UI.Songlist.ListChanged()
		})
	}
}

func (pms *PMS) handleEventPlayer() {
	pms.UI.App.PostFunc(func() {
		pms.UI.Songlist.SetCurrentSong(pms.CurrentSong())
	})
}

func (pms *PMS) handleEventMessage(msg message.Message) {
	message.Log(msg)
	pms.UI.App.PostFunc(func() {
		pms.UI.Multibar.SetMessage(msg)
	})
}
