package widgets

import (
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/version"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type UI struct {
	// UI elements
	App      *views.Application
	Layout   *views.BoxLayout
	Topbar   *views.TextBar
	Multibar *MultibarWidget
	Songlist *SongListWidget

	// Data resources
	Index           *index.Index
	defaultSongList *songlist.SongList

	views.WidgetWatchers
}

func NewUI() *UI {
	ui := &UI{}

	ui.App = &views.Application{}

	ui.Layout = views.NewBoxLayout(views.Vertical)
	ui.Topbar = views.NewTextBar()
	ui.Multibar = NewMultibarWidget()
	ui.Songlist = NewSongListWidget()

	ui.Layout.AddWidget(ui.Topbar, 0)
	ui.Layout.AddWidget(ui.Songlist, 2)
	ui.Layout.AddWidget(ui.Multibar, 0)

	ui.Multibar.Watch(ui)
	ui.Songlist.Watch(ui)

	str := version.ShortName() + " " + version.Version()
	style := tcell.StyleDefault
	style = style.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
	ui.Topbar.SetStyle(style)
	ui.Topbar.SetLeft(str, style)

	ui.Multibar.SetDefaultText("Type to search.")

	ui.App.SetRootWidget(ui)

	return ui
}

func (ui *UI) SetIndex(i *index.Index) {
	ui.Index = i
}

func (ui *UI) SetDefaultSonglist(s *songlist.SongList) {
	ui.defaultSongList = s
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
	ui.Layout.Resize()
}

func (ui *UI) SetView(v views.View) {
	ui.Layout.SetView(v)
}

func (ui *UI) Size() (int, int) {
	return ui.Layout.Size()
}

func (ui *UI) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {

	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlC:
			fallthrough
		case tcell.KeyCtrlD:
			ui.App.Quit()
			return true
		case tcell.KeyCtrlL:
			ui.App.Refresh()
			return true
		}

	case *EventListChanged:
		ui.App.Update()
		return true

	case *EventInputChanged:
		term := ui.Multibar.GetRuneString()
		ui.runIndexSearch(term)
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
	ui.Multibar.SetRight(str, tcell.StyleDefault)
}

func (ui *UI) runIndexSearch(term string) {
	if ui.Index == nil {
		return
	}
	if len(term) == 0 {
		ui.Songlist.SetCursor(0)
		ui.Songlist.SetSongList(ui.defaultSongList)
		return
	}
	if len(term) == 1 {
		return
	}
	results, err := ui.Index.Search(term)
	if err == nil {
		ui.Songlist.SetCursor(0)
		ui.Songlist.SetSongList(results)
		return
	}
}
