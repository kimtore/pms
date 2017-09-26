package pms

import (
	"time"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/message"
	"github.com/ambientsound/pms/term"
	"github.com/ambientsound/pms/topbar"
)

// Main does (eventually) read, evaluate, print, loop
func (pms *PMS) Main() {
	ticker := time.NewTicker(time.Millisecond * 1000)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pms.handleTicker()
		case <-pms.Connection.Connected:
			go pms.handleConnected()
		case subsystem := <-pms.Connection.IdleEvents:
			pms.handleEventIdle(subsystem)
		case <-pms.QuitSignal:
			pms.handleQuitSignal()
			return
		case key := <-pms.EventOption:
			pms.handleEventOption(key)
		case msg := <-pms.EventMessage:
			pms.handleEventMessage(msg)
		case e := <-pms.terminal.Events:
			pms.handleTerminalEvent(e)
		case s := <-pms.eventInputCommand:
			pms.Execute(s)
		}
		pms.ui.Draw()
	}
}

func (pms *PMS) handleTicker() {
	pms.database.SetPlayerStatus(pms.database.PlayerStatus().Tick())
}

func (pms *PMS) handleQuitSignal() {
	console.Log("Received quit signal, exiting.")
}

// handleTerminalEvent receives key input signals, checks the sequencer for key
// bindings, and runs commands if key bindings are found.
func (pms *PMS) handleTerminalEvent(e term.Event) {
	console.Log("%+v", e)

	switch e.Type {
	case term.EventKey:
		break
	case term.EventResize:
		pms.ui.Resize()
		return
	default:
		return
	}

	matches := pms.Sequencer.KeyInput(e.Key)
	seqString := pms.Sequencer.String()
	statusText := seqString

	input := pms.Sequencer.Match()
	if !matches || input != nil {
		// Reset statusbar if there is either no match or a complete match.
		statusText = ""
	}

	pms.EventMessage <- message.Sequencef(statusText)

	if input == nil {
		return
	}

	console.Log("Input sequencer matches bind: '%s' -> '%s'", seqString, input.Command)
	pms.eventInputCommand <- input.Command
}

func (pms *PMS) handleEventOption(key string) {
	console.Log("Option '%s' has been changed", key)
	switch key {
	case "topbar":
		pms.setupTopbar()
	case "columns":
		// list changed, FIXME
	}
}

func (pms *PMS) handleEventMessage(msg message.Message) {
	message.Log(msg)
	pms.ui.Multibar.SetMessage(msg)
}

// handleEventIdle triggers actions based on IDLE events.
func (pms *PMS) handleEventIdle(subsystem string) {
	var err error

	console.Log("MPD says it has IDLE events on the following subsystem: %s", subsystem)

	switch subsystem {
	case "database":
		err = pms.SyncLibrary()
	case "playlist":
		err = pms.SyncQueue()
	case "player":
		err = pms.UpdatePlayerStatus()
		if err != nil {
			break
		}
		err = pms.UpdateCurrentSong()
	case "options":
		err = pms.UpdatePlayerStatus()
	case "mixer":
		err = pms.UpdatePlayerStatus()
	default:
		console.Log("Ignoring updates by subsystem %s", subsystem)
	}

	if err != nil {
		pms.Error("Lost sync with MPD; reconnecting.")
		pms.Connection.Close()
	}
}

func (pms *PMS) setupTopbar() {
	config := pms.Options.StringValue("topbar")
	matrix, err := topbar.Parse(pms.API(), config)
	if err == nil {
		pms.ui.Topbar.SetMatrix(matrix)
		pms.ui.Resize()
	} else {
		pms.Error("Error in topbar configuration: %s", err)
	}
}
