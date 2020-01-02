package widgets

import (
	"fmt"
	"strings"

	"github.com/ambientsound/pms/api"
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
	Columnheaders *ColumnheadersWidget
	Multibar      *Multibar

	// Input events
	EventInputCommand chan string
	EventKeyInput     chan *tcell.EventKey

	// Data resources
	api          api.API
	searchResult songlist.Songlist

	// TCell
	view views.View
	style.Styled
	views.WidgetWatchers
}

func NewUI(a api.API) (*UI, error) {
	var err error

	ui := &UI{}

	ui.Screen, err = tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	ui.EventInputCommand = make(chan string, 16)
	ui.EventKeyInput = make(chan *tcell.EventKey, 16)

	ui.App = &views.Application{}
	ui.api = a

	ui.Topbar = NewTopbar()
	ui.Columnheaders = NewColumnheadersWidget()
	// ui.Multibar = NewMultibarWidget(ui.api, ui.EventKeyInput)

	ui.Multibar.Watch(ui)

	// Set styles
	ui.SetStylesheet(ui.api.Styles())
	ui.Topbar.SetStylesheet(ui.api.Styles())
	ui.Columnheaders.SetStylesheet(ui.api.Styles())
	ui.Multibar.SetStylesheet(ui.api.Styles())

	ui.CreateLayout()
	ui.App.SetScreen(ui.Screen)
	ui.App.SetRootWidget(ui)

	return ui, nil
}

func (ui *UI) CreateLayout() {
	ui.Layout = views.NewBoxLayout(views.Vertical)
	ui.Layout.AddWidget(ui.Topbar, 1)
	ui.Layout.AddWidget(ui.Columnheaders, 0)
	ui.Layout.AddWidget(ui.Multibar, 0)
	ui.Layout.SetView(ui.view)
}

func (ui *UI) Refresh() {
	ui.App.Refresh()
}

func (ui *UI) CurrentSonglistWidget() api.TableWidget {
	return nil
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
	ui.UpdateCursor()
	ui.Layout.Draw()
}

func (ui *UI) Resize() {
	ui.api.Db().Left().SetUpdated()
	ui.api.Db().Right().SetUpdated()
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

func (ui *UI) UpdateCursor() {
	/*
	switch ui.Multibar.multibar.Mode() {
	case multibar.ModeInput, multibar.ModeSearch:
		_, ymax := ui.Screen.Size()
		ui.Screen.ShowCursor(ui.Multibar.multibar.Cursor()+1, ymax-1)
	default:
		ui.Screen.HideCursor()
	}
	*/
}

func (ui *UI) PostFunc(f func()) {
	ui.App.PostFunc(f)
}

func (ui *UI) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {

	// If a list was changed, make sure we obtain the correct column widths.
	case *EventListChanged:
		tags := strings.Split(ui.api.Options().GetString(options.Columns), ",")
		cols := ui.api.Songlist().Columns(tags)
		ui.Columnheaders.SetColumns(cols)
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
	// str := ui.Songlist.PositionReadout()
	// ui.Multibar.SetRight(str, ui.Style("readout"))
}

func (ui *UI) runIndexSearch(term string) error {
	var err error

	library := ui.api.Library()
	if library == nil {
		return fmt.Errorf("Song library is not present.")
	}

	ui.searchResult, err = library.Search(term)

	ui.showSearchResult()

	return err
}

func (ui *UI) showSearchResult() {
	panel := ui.api.Db().Panel()
	if ui.searchResult != nil {
		panel.Activate(ui.searchResult)
	} else if panel.Last() != nil {
		panel.Activate(panel.Last())
	} else {
		panel.ActivateIndex(0)
	}
}
