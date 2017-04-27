package pms

import "github.com/ambientsound/pms/console"

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
		case <-pms.EventPlayer:
			pms.handleEventPlayer()
		case s := <-pms.EventMessage:
			pms.handleEventMessage(s)
		case s := <-pms.EventError:
			pms.handleEventError(s)
		case ev := <-pms.UI.EventKeyInput:
			pms.KeyInput(ev)
		case s := <-pms.UI.EventInputCommand:
			pms.Execute(s)
		}
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
		pms.UI.App.Update()
	})
}

func (pms *PMS) handleEventQueue() {
	console.Log("Queue updated in MPD, assigning to UI")
	pms.UI.App.PostFunc(func() {
		pms.UI.Songlist.ReplaceSonglist(pms.Queue)
		pms.UI.App.Update()
	})
}

func (pms *PMS) handleEventIndex() {
	console.Log("Search index updated, assigning to UI")
	pms.UI.App.PostFunc(func() {
		pms.UI.SetIndex(pms.Index)
	})
}

func (pms *PMS) handleEventPlayer() {
	pms.UI.App.PostFunc(func() {
		pms.UI.Playbar.SetPlayerStatus(pms.CurrentPlayerStatus())
		pms.UI.Playbar.SetSong(pms.CurrentSong())
		pms.UI.Songlist.SetCurrentSong(pms.CurrentSong())
		pms.UI.App.Update()
	})
}

func (pms *PMS) handleEventMessage(s string) {
	console.Log(s)
	pms.UI.App.PostFunc(func() {
		pms.UI.Multibar.SetText(s)
		pms.UI.App.Update()
	})
}

func (pms *PMS) handleEventError(s string) {
	console.Log(s)
	pms.UI.App.PostFunc(func() {
		pms.UI.Multibar.SetErrorText(s)
		pms.UI.App.Update()
	})
}
