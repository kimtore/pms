package prog

import (
	"bufio"
	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/log"
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/spotify/aggregator"
	"github.com/ambientsound/pms/spotify/auth"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/tabcomplete"
	"github.com/ambientsound/pms/widgets"
	"github.com/gdamore/tcell"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify"
	"io"
	"os"
	"strings"
)

type Visp struct {
	Auth   *spotify_auth.Handler
	Termui *widgets.Application

	client      spotify.Client
	commands    chan string
	clipboard   *songlist.BaseSonglist
	interpreter *input.Interpreter
	multibar    *multibar.Multibar
	options     *options.Options
	quit        chan interface{}
	sequencer   *keys.Sequencer
	stylesheet  style.Stylesheet
}

var _ api.API = &Visp{}

func (v *Visp) Clipboard() songlist.Songlist {
	return v.clipboard
}

func (v *Visp) Db() *db.Instance {
	return nil // FIXME
}

func (v *Visp) Exec(command string) error {
	log.Debugf("Run command: %s", command)
	return v.interpreter.Exec(command)
}

func (v *Visp) Library() *songlist.Library {
	return nil // FIXME
}

func (v *Visp) ListChanged() {
	// FIXME
}

func (v *Visp) Message(fmt string, a ...interface{}) {
	log.Infof(fmt, a...)
	log.Debugf("Using obsolete Message() for previous message")
}

func (v *Visp) MpdClient() *mpd.Client {
	log.Debugf("nil mpd client; might break")
	return nil // FIXME
}

func (v *Visp) OptionChanged(key string) {
	// FIXME
}

func (v *Visp) Options() *options.Options {
	return v.options
}

func (v *Visp) PlayerStatus() (p pms_mpd.PlayerStatus) {
	return // FIXME
}

func (v *Visp) Queue() *songlist.Queue {
	log.Debugf("nil queue; might break")
	return nil // FIXME
}

func (v *Visp) Quit() {
	v.quit <- new(interface{})
}

func (v *Visp) Sequencer() *keys.Sequencer {
	return v.sequencer
}

func (v *Visp) Multibar() *multibar.Multibar {
	return v.multibar
}

func (v *Visp) Song() *song.Song {
	log.Debugf("nil song; might break")
	return nil
}

func (v *Visp) Songlist() songlist.Songlist {
	log.Debugf("nil songlist; might break")
	return nil
}

func (v *Visp) Songlists() []songlist.Songlist {
	log.Debugf("nil songlists; might break")
	return nil // FIXME
}

func (v *Visp) Styles() style.Stylesheet {
	return v.stylesheet
}

func (v *Visp) UI() api.UI {
	return v.Termui
}

func (v *Visp) Init() {
	tcf := func(in string) multibar.TabCompleter {
		return tabcomplete.New(in, v)
	}
	v.commands = make(chan string, 1024)
	v.interpreter = input.NewCLI(v)
	v.multibar = multibar.New(tcf)
	v.options = options.New()
	v.quit = make(chan interface{}, 1)
	v.sequencer = keys.NewSequencer()
	v.stylesheet = make(style.Stylesheet)
}

func (v *Visp) Main() error {
	for {
		select {
		case token := <-v.Auth.Tokens():
			log.Debugf("Received Spotify token.")
			v.client = v.Auth.Client(token)
			viper.Set("spotify.accesstoken", token.AccessToken)
			viper.Set("spotify.refreshtoken", token.RefreshToken)
			err := viper.WriteConfig()
			if err != nil {
				log.Errorf("Unable to write configuration file: %s", err)
			}

		case <-v.quit:
			log.Infof("Exiting.")
			return nil

		// Send commands from the multibar into the main command queue.
		case command := <-v.multibar.Commands():
			v.commands <- command

		// Search input box. Discard for now.
		case query := <-v.multibar.Searches():
			lst, err := spotify_aggregator.Search(v.client, query)
			if err != nil {
				log.Errorf("spotify search: %s", err)
				break
			}
			v.UI().TableWidget().SetList(lst)
			v.UI().TableWidget().SetColumns(lst.ColumnNames())

		// Process the command queue.
		case command := <-v.commands:
			err := v.Exec(command)
			if err != nil {
				log.Errorf(err.Error())
				v.multibar.Error(err)
			}

		// Try handling the input event in the multibar.
		// If multibar is disabled (input mode = normal), try handling the event in the UI layer.
		// If unhandled still, run it through the keyboard binding maps to try to get a command.
		case ev := <-v.Termui.Events():
			if v.multibar.Input(ev) {
				break
			}
			if v.Termui.HandleEvent(ev) {
				break
			}
			cmd := v.keyEventCommand(ev)
			if len(cmd) == 0 {
				break
			}
			v.commands <- cmd
		}

		// Draw UI after processing any event.
		v.Termui.Draw()
	}
}

// KeyInput receives key input signals, checks the sequencer for key bindings,
// and runs commands if key bindings are found.
func (v *Visp) keyEventCommand(event tcell.Event) string {
	ev, ok := event.(*tcell.EventKey)
	if !ok {
		return ""
	}

	matches := v.sequencer.KeyInput(ev)
	seqString := v.sequencer.String()

	match := v.sequencer.Match()
	if !matches || match != nil {
		// Reset statusbar if there is either no match or a complete match.
		seqString = ""
	}

	log.Debugf("console: %s", seqString)

	if match == nil {
		return ""
	}

	log.Debugf("Input sequencer matches bind: '%s' -> '%s'", seqString, match.Command)

	return match.Command
}

// SourceDefaultConfig reads, parses, and executes the default config.
func (v *Visp) SourceDefaultConfig() error {
	reader := strings.NewReader(options.Defaults)
	return v.SourceConfig(reader)
}

// SourceConfigFile reads, parses, and executes a config file.
func (v *Visp) SourceConfigFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return v.SourceConfig(file)
}

// SourceConfig reads, parses, and executes config lines.
func (v *Visp) SourceConfig(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		err := v.interpreter.Execute(scanner.Text())
		if err != nil {
			return err
		}
	}
	return nil
}
