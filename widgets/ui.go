package widgets

import (
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/songlist"
	termbox "github.com/nsf/termbox-go"

	"github.com/gdamore/tcell"
)

type UI struct {
	// UI elements
	Topbar        *Topbar
	Columnheaders *ColumnheadersWidget
	Multibar      *MultibarWidget
	Songlist      *SonglistWidget

	// Input events
	EventInputCommand chan string
	EventKeyInput     chan *tcell.EventKey

	// Data resources
	api          api.API
	searchResult songlist.Songlist
}

func NewUI(a api.API) *UI {
	var err error

	ui := &UI{}

	err = termbox.Init()
	if err != nil {
		return nil
	}

	termbox.HideCursor()

	ui.EventInputCommand = make(chan string, 16)
	ui.EventKeyInput = make(chan *tcell.EventKey, 16)

	ui.api = a

	// Initialize widgets
	ui.Topbar = NewTopbar()
	ui.Columnheaders = NewColumnheadersWidget()
	ui.Multibar = NewMultibarWidget(ui.api, ui.EventKeyInput)
	ui.Songlist = NewSonglistWidget(ui.api)

	// Set styles
	ui.Topbar.SetStylesheet(ui.api.Styles())
	ui.Columnheaders.SetStylesheet(ui.api.Styles())
	ui.Songlist.SetStylesheet(ui.api.Styles())
	ui.Multibar.SetStylesheet(ui.api.Styles())

	return ui
}

func (ui *UI) Refresh() {
}

func (ui *UI) CurrentSonglistWidget() api.SonglistWidget {
	return ui.Songlist
}

func (ui *UI) Resize() {
	ui.api.Db().Left().SetUpdated()
	ui.api.Db().Right().SetUpdated()
	/*
		ui.CreateLayout()
		ui.Layout.Resize()
		ui.PostEventWidgetResize(ui)
	*/
}

func (ui *UI) UpdateCursor() {
	/*
		switch ui.Multibar.Mode() {
		case constants.MultibarModeInput, constants.MultibarModeSearch:
			_, ymax := ui.Screen.Size()
			ui.Screen.ShowCursor(ui.Multibar.Cursor()+1, ymax-1)
		default:
			ui.Screen.HideCursor()
		}
	*/
}

/*
func (ui *UI) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {

	// If a list was changed, make sure we obtain the correct column widths.
	case *EventListChanged:
		tags := strings.Split(ui.options.StringValue("columns"), ",")
		cols := ui.api.Songlist().Columns(tags)
		ui.Songlist.SetColumns(tags)
		ui.Columnheaders.SetColumns(cols)
		return true

	case *EventInputChanged:
		term := ui.Multibar.RuneString()
		mode := ui.Multibar.Mode()
		switch mode {
		case constants.MultibarModeSearch:
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
		case constants.MultibarModeInput:
			ui.EventInputCommand <- term
		case constants.MultibarModeSearch:
			if ui.searchResult != nil {
				if ui.searchResult.Len() > 0 {
					ui.api.Db().Panel().Add(ui.searchResult)
				} else {
					ui.searchResult = nil
				}
			}
			ui.showSearchResult()
		}
		ui.Multibar.SetMode(constants.MultibarModeNormal)
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
*/
