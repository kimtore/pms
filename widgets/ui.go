package widgets

import (
	"fmt"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/input/parser"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/version"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type StyleMap map[string]tcell.Style

type UI struct {
	// UI elements
	App    *views.Application
	Layout *views.BoxLayout

	Topbar        *views.TextBar
	Playbar       *PlaybarWidget
	Columnheaders *ColumnheadersWidget
	Multibar      *MultibarWidget
	Songlist      *SonglistWidget
	Styles        StyleMap

	// Input events
	EventInputCommand chan string
	EventKeyInput     chan parser.KeyEvent

	// Data resources
	Index        *index.Index
	options      *options.Options
	searchResult songlist.Songlist

	// TCell
	view views.View
	widget
	views.WidgetWatchers
}

func NewUI(opts *options.Options) *UI {
	ui := &UI{}

	ui.EventInputCommand = make(chan string, 16)
	ui.EventKeyInput = make(chan parser.KeyEvent, 16)

	ui.App = &views.Application{}
	ui.options = opts

	ui.Topbar = views.NewTextBar()
	ui.Playbar = NewPlaybarWidget()
	ui.Columnheaders = NewColumnheadersWidget()
	ui.Multibar = NewMultibarWidget(ui.EventKeyInput)
	ui.Songlist = NewSonglistWidget(ui.options)

	ui.Multibar.Watch(ui)
	ui.Songlist.Watch(ui)
	ui.Playbar.Watch(ui)

	// Set styles
	ui.Styles = make(StyleMap)
	ui.SetStyleMap(ui.Styles)
	ui.Columnheaders.SetStyleMap(ui.Styles)
	ui.Playbar.SetStyleMap(ui.Styles)
	ui.Songlist.SetStyleMap(ui.Styles)
	ui.Multibar.SetStyleMap(ui.Styles)

	// Styles for widgets that don't have their own class yet.
	ui.Topbar.SetStyle(ui.Style("topbar"))
	ui.Topbar.SetLeft(version.ShortName(), ui.Style("topbar"))
	ui.Topbar.SetRight(version.Version(), ui.Style("topbar"))

	ui.CreateLayout()
	ui.App.SetRootWidget(ui)

	return ui
}

func (ui *UI) CreateLayout() {
	ui.Layout = views.NewBoxLayout(views.Vertical)
	ui.Layout.AddWidget(ui.Topbar, 0)
	ui.Layout.AddWidget(ui.Playbar, 0)
	ui.Layout.AddWidget(ui.Columnheaders, 0)
	ui.Layout.AddWidget(ui.Songlist, 2)
	ui.Layout.AddWidget(ui.Multibar, 0)
	ui.Layout.SetView(ui.view)
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

func (ui *UI) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {

	case *EventListChanged:
		ui.Topbar.SetCenter(ui.Title(), ui.Style("title"))
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
