package widgets

import (
	"github.com/ambientsound/pms/pms"
	"github.com/ambientsound/pms/utils"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type PlaybarWidget struct {
	status pms.PlayerStatus
	view   views.View

	views.WidgetWatchers
}

var playRunes = map[string]rune{
	pms.StatePlay:    '▶',
	pms.StatePause:   '⏸',
	pms.StateStop:    '⏹',
	pms.StateUnknown: '�',
}

func NewPlaybarWidget() *PlaybarWidget {
	return &PlaybarWidget{}
}

func (w *PlaybarWidget) SetPlayerStatus(s pms.PlayerStatus) {
	w.status = s
	w.PostEventWidgetContent(w)
}

func (w *PlaybarWidget) drawNext(x, y int, s string, style tcell.Style) int {
	p := 0
	for p = 0; p < len(s); p++ {
		w.view.SetContent(x+p, y, rune(s[p]), nil, style)
	}
	return x + p
}

func (w *PlaybarWidget) Draw() {
	x, y := 0, 0
	style := tcell.StyleDefault

	x = w.drawNext(x, y, string(playRunes[w.status.State]), style)
	x = w.drawNext(x+1, y, w.status.State, style)
	x = w.drawNext(x+1, y, utils.TimeString(int(w.status.Elapsed)), style)
	x = w.drawNext(x+1, y, "/", style)
	x = w.drawNext(x+1, y, utils.TimeString(w.status.Time), style)
}

func (w *PlaybarWidget) SetView(v views.View) {
	w.view = v
}

func (w *PlaybarWidget) Size() (int, int) {
	x, y := w.view.Size()
	y = 1
	return x, y
}

func (w *PlaybarWidget) Resize() {
}

func (w *PlaybarWidget) HandleEvent(ev tcell.Event) bool {
	return false
}
