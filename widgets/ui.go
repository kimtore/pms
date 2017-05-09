package widgets

import (
	"fmt"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/style"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type UI struct {
	// UI elements
	Screen tcell.Screen
	App    *views.Application
	Layout *views.BoxLayout

	Topbar        *Topbar
	Playbar       *PlaybarWidget
	Columnheaders *ColumnheadersWidget
	Multibar      *MultibarWidget
	Songlist      *SonglistWidget
	Stylesheet    style.Stylesheet

	// Input events
	EventInputCommand chan string
	EventKeyInput     chan parser.KeyEvent

	// Data resources
	Index        *index.Index
	options      *options.Options
	searchResult songlist.Songlist

	// TCell
	view views.View
	style.Styled
	views.WidgetWatchers
}

func NewUI(opts *options.Options) *UI {
	var err error

	ui := &UI{}

	ui.Screen, err = tcell.NewScreen()
	if err != nil {
		return nil
	}

	ui.EventInputCommand = make(chan string, 16)
	ui.EventKeyInput = make(chan parser.KeyEvent, 16)

	ui.App = &views.Application{}
	ui.options = opts

	ui.Topbar = NewTopbar()
	ui.Playbar = NewPlaybarWidget()
	ui.Columnheaders = NewColumnheadersWidget()
	ui.Multibar = NewMultibarWidget(ui.EventKeyInput)
	ui.Songlist = NewSonglistWidget(ui.options)

	ui.Multibar.Watch(ui)
	ui.Songlist.Watch(ui)
	ui.Playbar.Watch(ui)

	// Set styles
	ui.Stylesheet = make(style.Stylesheet)
	ui.SetStylesheet(ui.Stylesheet)
	ui.Topbar.SetStylesheet(ui.Stylesheet)
	ui.Columnheaders.SetStylesheet(ui.Stylesheet)
	ui.Playbar.SetStylesheet(ui.Stylesheet)
	ui.Songlist.SetStylesheet(ui.Stylesheet)
	ui.Multibar.SetStylesheet(ui.Stylesheet)

	ui.CreateLayout()
	ui.App.SetScreen(ui.Screen)
	ui.App.SetRootWidget(ui)

	return ui
}

func (ui *UI) CreateLayout() {
	ui.Layout = views.NewBoxLayout(views.Vertical)
	ui.Layout.AddWidget(ui.Topbar, 1)
	ui.Layout.AddWidget(ui.Playbar, 0)
	ui.Layout.AddWidget(ui.Columnheaders, 0)
	ui.Layout.AddWidget(ui.Songlist, 2)
	ui.Layout.AddWidget(ui.Multibar, 0)
	ui.Layout.SetView(ui.view)
}

func (ui *UI) Refresh() {
	ui.App.Refresh()
}

func (ui *UI) SetIndex(i *index.Index) {
	ui.Index = i
}

func (ui *UI) CurrentSonglistWidget() *SonglistWidget {
	return ui.Songlist
}

func (ui *UI) Start() {
	ui.App.Start()
}

func (ui *UI) Wait() error {
	return ui.App.Wait()
}

func (ui *UI) Quit() {
	ui.App.Quit()
}

func (ui *UI) Draw() {
	ui.Layout.Draw()
}

func (ui *UI) Resize() {
	ui.CreateLayout()
	ui.Layout.Resize()
	ui.PostEventWidgetResize(ui)
}

func (ui *UI) SetView(v views.View) {
	ui.view = v
	ui.Layout.SetView(v)
}

func (ui *UI) Size() (int, int) {
	return ui.view.Size()
}

func (ui *UI) Title() string {
	index, err := ui.Songlist.SonglistIndex()
	if err == nil {
		return fmt.Sprintf("[%d/%d] %s", index+1, ui.Songlist.SonglistsLen(), ui.Songlist.Name())
	} else {
		return fmt.Sprintf("[...] %s", ui.Songlist.Name())
	}
}

func (ui *UI) UpdateCursor() {
	switch ui.Multibar.Mode() {
	case MultibarModeInput, MultibarModeSearch:
		_, ymax := ui.Screen.Size()
		ui.Screen.ShowCursor(ui.Multibar.RuneLen()+1, ymax-1)
	default:
		ui.Screen.HideCursor()
	}
}

func (ui *UI) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {

	case *EventListChanged:
		//ui.Topbar.SetCenter(ui.Title(), ui.Style("title"))
		ui.Columnheaders.SetColumns(ui.Songlist.Columns())
		ui.App.Update()
		return true

	case *EventModeSync:
		console.Log("EventModeChanged %d", ev.InputMode)
		switch {
		case ev.InputMode != ui.Multibar.Mode():
			console.Log("Resetting multibar mode based on songlist change")
			ui.Multibar.SetMode(ev.InputMode)
		case ev.InputMode == MultibarModeVisual && !ui.Songlist.HasVisualSelection():
			console.Log("Enabling visual selection based on multibar setting")
			ui.Songlist.EnableVisualSelection()
		case ev.InputMode != MultibarModeVisual && ui.Songlist.HasVisualSelection():
			console.Log("Disabling visual selection based on multibar setting")
			ui.Songlist.DisableVisualSelection()
		}
		ui.UpdateCursor()
		return true

	case *EventInputChanged:
		term := ui.Multibar.RuneString()
		mode := ui.Multibar.Mode()
		switch mode {
		case MultibarModeSearch:
			if err := ui.runIndexSearch(term); err != nil {
				console.Log("Error while searching: %s", err)
			}
		}
		ui.UpdateCursor()
		return true

	case *EventInputFinished:
		term := ui.Multibar.RuneString()
		mode := ui.Multibar.Mode()
		switch mode {
		case MultibarModeInput:
			ui.EventInputCommand <- term
		case MultibarModeSearch:
			if ui.searchResult != nil {
				if ui.searchResult.Len() > 0 {
					ui.Songlist.AddSonglist(ui.searchResult)
				} else {
					ui.searchResult = nil
				}
			}
			ui.showSearchResult()
		}
		ui.Multibar.SetMode(MultibarModeNormal)
		return true

	case *EventScroll:
		ui.refreshPositionReadout()
		return true
	}

	if ui.Layout.HandleEvent(ev) {
		return true
	}

	return false
}

func (ui *UI) refreshPositionReadout() {
	str := ui.Songlist.PositionReadout()
	ui.Multibar.SetRight(str, ui.Style("readout"))
}

func (ui *UI) runIndexSearch(term string) error {
	var err error

	if ui.Index == nil {
		return fmt.Errorf("Search index is not operational")
	}

	ui.searchResult, err = ui.Index.Search(term)

	ui.showSearchResult()

	return err
}

func (ui *UI) showSearchResult() {
	if ui.searchResult != nil {
		ui.Songlist.SetSonglist(ui.searchResult)
	} else if ui.Songlist.FallbackSonglist() != nil {
		ui.Songlist.SetSonglist(ui.Songlist.FallbackSonglist())
	} else {
		ui.Songlist.SetSonglistIndex(0)
	}
}
