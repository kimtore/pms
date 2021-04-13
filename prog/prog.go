package prog

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ambientsound/pms/api"
	"github.com/ambientsound/pms/commands"
	"github.com/ambientsound/pms/db"
	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/list"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/multibar"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/player"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/spotify/aggregator"
	"github.com/ambientsound/pms/spotify/library"
	spotify_proxyclient "github.com/ambientsound/pms/spotify/proxyclient"
	"github.com/ambientsound/pms/style"
	"github.com/ambientsound/pms/tabcomplete"
	"github.com/ambientsound/pms/tokencache"
	"github.com/ambientsound/pms/widgets"
	"github.com/gdamore/tcell"
	"github.com/google/uuid"
	"github.com/zmb3/spotify"
)

const (
	refreshTokenTimeout       = time.Second * 5
	refreshTokenRetryInterval = time.Second * 30
)

type Visp struct {
	Termui     *widgets.Application
	Tokencache tokencache.Tokencache

	client       *spotify.Client
	clipboard    *songlist.BaseSonglist
	commands     chan string
	db           *db.List
	interpreter  *input.Interpreter
	library      *spotify_library.List
	list         list.List
	multibar     *multibar.Multibar
	player       player.State
	quit         chan interface{}
	sequencer    *keys.Sequencer
	stylesheet   style.Stylesheet
	ticker       *time.Ticker
	tokenRefresh <-chan time.Time
}

var _ api.API = &Visp{}

func (v *Visp) Init() {
	tcf := func(in string) multibar.TabCompleter {
		return tabcomplete.New(in, v)
	}
	v.db = db.New()
	v.commands = make(chan string, 1024)
	v.interpreter = input.NewCLI(v)
	v.library = spotify_library.New()
	v.multibar = multibar.New(tcf)
	v.quit = make(chan interface{}, 1)
	v.sequencer = keys.NewSequencer()
	v.stylesheet = make(style.Stylesheet)
	v.ticker = time.NewTicker(time.Second)
	v.tokenRefresh = make(chan time.Time)

	v.SetList(log.List(log.InfoLevel))
}

func (v *Visp) Main() error {
	for {
		select {
		case <-v.quit:
			log.Infof("Exiting.")
			return nil

		case <-v.ticker.C:
			err := v.updatePlayer()
			if err != nil {
				log.Errorf("update player: %s", err)
			}

		case <-v.tokenRefresh:
			log.Infof("Spotify access token is too old, refreshing...")
			err := v.refreshToken()
			if err != nil {
				log.Errorf("Refresh Spotify access token: %s", err)
			}

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
			lst.SetID(uuid.New().String())
			lst.SetName(fmt.Sprintf("Search for '%s'", query))
			lst.SetVisibleColumns(strings.Split(columns, ","))
			v.SetList(lst)

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

func (v *Visp) updatePlayer() error {
	var err error

	now := time.Now()
	pollInterval := time.Second * time.Duration(v.Options().GetInt(options.PollInterval))

	// no time for polling yet; just increase the ticker.
	if v.player.CreateTime.Add(pollInterval).After(now) {
		v.player = v.player.Tick()
		return nil
	}

	log.Debugf("fetching new player information")

	client, err := v.Spotify()
	if err != nil {
		return err
	}

	state, err := client.PlayerState()
	if err != nil {
		return err
	}

	v.player = player.NewState(*state)

	return nil
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

func (v *Visp) refreshToken() error {
	server := v.Options().GetString(options.SpotifyAuthServer)
	client := &http.Client{
		Timeout: refreshTokenTimeout,
	}
	token, err := spotify_proxyclient.RefreshToken(server, client, v.Tokencache.Cached())
	if err != nil {
		v.tokenRefresh = time.After(refreshTokenRetryInterval)
		return err
	}
	return v.Authenticate(token)
}
