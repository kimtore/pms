package prog

import (
	"bufio"
	"fmt"
	"github.com/ambientsound/gompd/mpd"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/log"
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/spotify/aggregator"
	"github.com/ambientsound/pms/spotify/auth"
	"github.com/ambientsound/pms/spotify/tracklist"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/tabcomplete"
	"github.com/ambientsound/pms/tokencache"
	"github.com/ambientsound/pms/topbar"
	"github.com/ambientsound/pms/widgets"
	"github.com/gdamore/tcell"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io"
	"os"
	"strings"
)

type Visp struct {
	Auth       *spotify_auth.Handler
	Termui     *widgets.Application
	Tokencache tokencache.Tokencache

	client      *spotify.Client
	commands    chan string
	clipboard   *songlist.BaseSonglist
	interpreter *input.Interpreter
	multibar    *multibar.Multibar
	quit        chan interface{}
	sequencer   *keys.Sequencer
	stylesheet  style.Stylesheet
}

var _ api.API = &Visp{}

func (v *Visp) Authenticate() error {
	err := v.SetupAuthenticator()
	if err != nil {
		return fmt.Errorf("cannot authenticate with Spotify: %s", err)
	}
	url := v.Auth.AuthURL()
	log.Infof("Please authenticate with Spotify at: %s", url)

	return nil
}

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

func (v *Visp) List() list.List {
	return v.Termui.TableWidget().List()
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
	switch key {
	case options.LogFile:
		logFile := v.Options().GetString(options.LogFile)
		overwrite := v.Options().GetBool(options.LogOverwrite)
		if len(logFile) == 0 {
			break
		}
		err := log.Configure(logFile, overwrite)
		if err != nil {
			log.Errorf("log configuration: %s", err)
			break
		}
		log.Infof("Note: log file will be backfilled with existing log")
		log.Infof("Writing debug log to %s", logFile)

	case options.Topbar:
		config := v.Options().GetString(options.Topbar)
		matrix, err := topbar.Parse(v, config)
		if err == nil {
			_ = matrix
			// FIXME
			// v.Termui.Widgets.Topbar.SetMatrix(matrix)
		} else {
			log.Errorf("topbar configuration: %s", err)
		}
	}
}

func (v *Visp) Options() api.Options {
	return viper.GetViper()
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

func (v *Visp) Spotify() (*spotify.Client, error) {
	if v.client == nil {
		return nil, fmt.Errorf("please run `auth` to authenticate with Spotify")
	}
	err := v.SetupAuthenticator()
	if err != nil {
		return nil, fmt.Errorf("unable to obtain Spotify client: %s", err.Error())
	}
	token, err := v.client.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to refresh Spotify token: %s", err)
	}
	_ = v.Tokencache.Write(*token)
	return v.client, nil
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

func (v *Visp) Tracklist() *spotify_tracklist.List {
	switch v := v.UI().TableWidget().List().(type) {
	case *spotify_tracklist.List:
		return v
	default:
		return nil
	}
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
	v.quit = make(chan interface{}, 1)
	v.sequencer = keys.NewSequencer()
	v.stylesheet = make(style.Stylesheet)
}

func (v *Visp) Main() error {
	for {
		select {
		case token := <-v.Auth.Tokens():
			v.SetToken(token)

		case <-v.quit:
			log.Infof("Exiting.")
			return nil

		// Send commands from the multibar into the main command queue.
		case command := <-v.multibar.Commands():
			v.commands <- command

		// Search input box.
		case query := <-v.multibar.Searches():
			client, err := v.Spotify()
			if err != nil {
				log.Errorf(err.Error())
				break
			}
			lst, err := spotify_aggregator.Search(*client, query, v.Options().GetInt(options.Limit))
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
	log.Infof("Reading configuration file %s", path)
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

func (v *Visp) SetupAuthenticator() error {
	clientID := v.Options().GetString(options.SpotifyClientID)
	clientSecret := v.Options().GetString(options.SpotifyClientSecret)

	if len(clientID) == 0 && len(clientSecret) == 0 {
		return fmt.Errorf("you must configure `%s` and `%s`", options.SpotifyClientID, options.SpotifyClientSecret)
	} else if len(clientID) == 0 {
		return fmt.Errorf("you must configure `%s`", options.SpotifyClientID)
	} else if len(clientSecret) == 0 {
		return fmt.Errorf("you must configure `%s`", options.SpotifyClientSecret)
	}

	v.Auth.SetCredentials(clientID, clientSecret)

	return nil
}

func (v *Visp) SetToken(token oauth2.Token) {
	log.Infof("Received Spotify access token.")
	cli := v.Auth.Client(token)
	v.client = &cli
	err := v.Tokencache.Write(token)
	if err != nil {
		log.Errorf("Unable to write Spotify token to file: %s", err)
	}
}
