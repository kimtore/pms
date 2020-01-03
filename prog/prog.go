package prog

import (
	"bufio"
	"fmt"
	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/spotify/aggregator"
	"github.com/ambientsound/pms/spotify/auth"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/tabcomplete"
	"github.com/ambientsound/pms/tokencache"
	"github.com/ambientsound/pms/widgets"
	"github.com/gdamore/tcell"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io"
	"os"
	"strings"
	"time"
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
	ticker      *time.Ticker
}

var _ api.API = &Visp{}

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
	v.ticker = time.NewTicker(time.Hour)
	v.ticker.Stop()
}

func (v *Visp) Main() error {
	for {
		select {
		case token := <-v.Auth.Tokens():
			v.SetToken(token)

		case <-v.quit:
			log.Infof("Exiting.")
			return nil

		case <-v.ticker.C:
			// poll spotify

		// Send commands from the multibar into the main command queue.
		case command := <-v.multibar.Commands():
			v.commands <- command

		// Search input box.
		case query := <-v.multibar.Searches():
			if len(query) == 0 {
				break
			}
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
			columns := v.Options().GetString(options.Columns)
			sort := v.Options().GetString(options.Sort)
			err = lst.Sort(strings.Split(sort, ","))
			if err != nil {
				log.Errorf("sort search results: %s")
			}
			v.UI().TableWidget().SetList(lst)
			v.UI().TableWidget().SetColumns(strings.Split(columns, ","))

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

	contexts := commands.Contexts(v)
	v.sequencer.KeyInput(ev, contexts)
	match := v.sequencer.Match(contexts)

	if match == nil {
		return ""
	}

	log.Debugf("Input sequencer matches bind: '%s' -> '%s'", match.Sequence, match.Command)

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
		err := v.interpreter.Exec(scanner.Text())
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Visp) setupAuthenticator() error {
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
