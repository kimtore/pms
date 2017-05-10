package widgets

import (
	"fmt"

	"github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/utils"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type PlaybarWidget struct {
	status mpd.PlayerStatus
	view   views.View
	song   song.Song

	style.Styled
	views.WidgetWatchers
}

var playRunes = map[string]rune{
	mpd.StatePlay:    '\u25b6',
	mpd.StatePause:   '\u23f8',
	mpd.StateStop:    '\u23f9',
	mpd.StateUnknown: '\u2bd1',
}

var playStrings = map[string]string{
	mpd.StatePlay:    "|>",
	mpd.StatePause:   "||",
	mpd.StateStop:    "[]",
	mpd.StateUnknown: "??",
}

func StatusRune(r rune, val bool) rune {
	if val {
		return r
	}
	return '-'
}

func NewPlaybarWidget() *PlaybarWidget {
	return &PlaybarWidget{}
}

func (w *PlaybarWidget) SetPlayerStatus(s mpd.PlayerStatus) {
	w.status = s
	w.PostEventWidgetContent(w)
}

func (w *PlaybarWidget) SetSong(s *song.Song) {
	if s != nil {
		w.song = *s
	} else {
		w.song = song.Song{}
	}
	w.PostEventWidgetContent(w)
}

func (w *PlaybarWidget) drawNext(x, y int, runes []rune, style tcell.Style) int {
	strlen := 0
	for p, r := range runes {
		w.view.SetContent(x+p, y, r, nil, style)
		strlen++
	}
	return x + strlen
}

func (w *PlaybarWidget) drawNextChar(x, y int, r rune, style tcell.Style) int {
	w.view.SetContent(x, y, r, nil, style)
	return x + 1
}

func (w *PlaybarWidget) Draw() {
	x, y := 0, 0

	// 54% ----   00:00 ■ 00:00   Artist - Title

	volume := fmt.Sprintf("%d%%", w.status.Volume)
	x = w.drawNext(x, y, []rune(volume), w.Style("volume"))

	x = w.drawNextChar(x+1, y, StatusRune('c', w.status.Consume), w.Style("switches"))
	x = w.drawNextChar(x+0, y, StatusRune('z', w.status.Random), w.Style("switches"))
	x = w.drawNextChar(x+0, y, StatusRune('s', w.status.Single), w.Style("switches"))
	x = w.drawNextChar(x+0, y, StatusRune('r', w.status.Repeat), w.Style("switches"))

	x = w.drawNext(x+1, y, []rune(utils.TimeString(int(w.status.Elapsed))), w.Style("elapsed"))
	x = w.drawNext(x+1, y, []rune(playStrings[w.status.State]), w.Style("symbol"))
	w.drawNext(x+1, y, []rune(utils.TimeString(w.status.Time)), w.Style("time"))
}

func (w *PlaybarWidget) SetView(v views.View) {
	w.view = v
}

func (w *PlaybarWidget) Size() (int, int) {
	x, _ := w.view.Size()
	return x, 2
}

func (w *PlaybarWidget) Resize() {
}

func (w *PlaybarWidget) HandleEvent(ev tcell.Event) bool {
	return false
}
