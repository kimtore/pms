package widgets

import (
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/term"
	termbox "github.com/nsf/termbox-go"
)

type UI struct {
	// UI elements
	Topbar   *Topbar
	Multibar *MultibarWidget
	Songlist *SonglistWidget

	// Data resources
	api          api.API
	searchResult songlist.Songlist
}

func NewUI(a api.API) (*UI, error) {
	var err error

	ui := &UI{}

	err = termbox.Init()
	if err != nil {
		return nil, err
	}

	termbox.HideCursor()
	termbox.SetOutputMode(termbox.Output256)

	ui.api = a

	// Initialize widgets
	ui.Topbar = NewTopbar()
	ui.Multibar = NewMultibarWidget(ui.api)
	ui.Songlist = NewSonglistWidget(ui.api)

	// Set styles
	ui.Topbar.SetStylesheet(ui.api.Styles())
	ui.Songlist.SetStylesheet(ui.api.Styles())
	ui.Multibar.SetStylesheet(ui.api.Styles())

	// Layout all widgets
	ui.Resize()

	return ui, nil
}

// Resize calculates the size of all widgets.
func (ui *UI) Resize() {
	// FIXME: set dirty
	//ui.api.Db().Left().SetUpdated()
	//ui.api.Db().Right().SetUpdated()
	termbox.Flush()
	w, h := termbox.Size()
	console.Log("Terminal resized: %dx%d characters", w, h)
	ui.Topbar.SetCanvas(term.NewCanvas(0, 0, w, ui.Topbar.Height()))
	ui.Songlist.SetCanvas(term.NewCanvas(0, ui.Topbar.Height(), w, h-ui.Topbar.Height()-1))
	ui.Songlist.Resize()
}

func (ui *UI) Draw() {
	timer := time.Now()

	//t := time.Now()
	ui.Topbar.Draw()
	//console.Log("Topbar::Draw() in %s", time.Since(t).String())

	//t = time.Now()
	ui.Songlist.Draw()
	//console.Log("Songlist::Draw() in %s", time.Since(t).String())

	//t = time.Now()
	termbox.Flush()
	//console.Log("Terminal flush in %s", time.Since(t).String())

	console.Log("UI::Draw() total duration %s", time.Since(timer).String())
}

func (ui *UI) CurrentSonglistWidget() api.SonglistWidget {
	return ui.Songlist
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
