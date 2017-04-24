package pms

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/input"
	"github.com/ambientsound/pms/input/keys"
	"github.com/ambientsound/pms/input/parser"
	pms_mpd "github.com/ambientsound/pms/mpd"
	"github.com/ambientsound/pms/options"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/widgets"
	"github.com/ambientsound/pms/xdg"

	"github.com/fhs/gompd/mpd"
)

// PMS is a kitchen sink of different objects, glued together as a singleton class.
type PMS struct {
	MpdStatus        pms_mpd.PlayerStatus
	MpdClient        *mpd.Client
	MpdClientWatcher *mpd.Watcher
	currentSong      *song.Song
	Index            *index.Index
	CLI              *input.CLI
	UI               *widgets.UI
	Queue            *songlist.Queue
	Library          *songlist.Library
	Options          *options.Options
	Sequencer        *keys.Sequencer
	mutex            sync.Mutex

	ticker chan time.Time

	host     string
	port     string
	password string

	queueVersion   int
	libraryVersion int
	indexVersion   int

	EventError   chan string
	EventIndex   chan int
	EventLibrary chan int
	EventQueue   chan int
	EventMessage chan string
	EventPlayer  chan int
	QuitSignal   chan int
}

func createDirectory(dir string) error {
	dir_mode := os.ModeDir | 0755
	return os.MkdirAll(dir, dir_mode)
}

func makeAddress(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}

func indexDirectory(host, port string) string {
	cache_dir := xdg.CacheDirectory()
	index_dir := path.Join(cache_dir, host, port, "index")
	return index_dir
}

func indexStateFile(host, port string) string {
	cache_dir := xdg.CacheDirectory()
	state_file := path.Join(cache_dir, host, port, "state")
	return state_file
}

func (pms *PMS) writeIndexStateFile(version int) error {
	path := indexStateFile(pms.host, pms.port)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	str := fmt.Sprintf("%d\n", version)
	file.WriteString(str)
	return nil
}

func (pms *PMS) readIndexStateFile() (int, error) {
	path := indexStateFile(pms.host, pms.port)
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		version, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, err
		}
		return version, nil
	}

	return 0, fmt.Errorf("No data in index file")
}

func (pms *PMS) Message(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	pms.EventMessage <- s
}

func (pms *PMS) Error(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	pms.EventError <- s
}

func (pms *PMS) SetConnectionParams(host, port, password string) {
	pms.MpdClient = nil
	pms.host = host
	pms.port = port
	pms.password = password
}

func (pms *PMS) LoopConnect() {
	for {
		err := pms.Connect()
		if err == nil {
			return
		}
		pms.Error("Error while connecting to MPD: %s", err)
		time.Sleep(1 * time.Second)
	}
}

func (pms *PMS) Connect() error {
	var err error

	addr := makeAddress(pms.host, pms.port)

	pms.MpdClient = nil
	pms.MpdClientWatcher = nil

	pms.Message("Establishing MPD IDLE connection to %s...", addr)

	pms.MpdClientWatcher, err = mpd.NewWatcher(`tcp`, addr, pms.password)
	if err != nil {
		pms.Error("Connection error: %s", err)
		goto errors
	}
	pms.Message("Connected to %s.", addr)

	err = pms.PingConnect()
	if err != nil {
		goto errors
	}

	go pms.watchMpdIdleErrors()
	go pms.watchMpdIdleEvents()
	go pms.runTicker()

	/*
		err = pms.UpdatePlayerStatus()
		if err != nil {
			goto errors
		}
	*/

	err = pms.UpdateCurrentSong()
	if err != nil {
		goto errors
	}

	err = pms.SyncQueue()
	if err != nil {
		goto errors
	}

	err = pms.SyncLibrary()
	if err != nil {
		goto errors
	}

	pms.Message("Ready.")

	return nil

errors:

	pms.Error("ERROR: %s", err)

	if pms.MpdClient != nil {
		pms.MpdClient.Close()
	}
	if pms.MpdClientWatcher != nil {
		pms.MpdClientWatcher.Close()
	}
	if pms.ticker != nil {
		close(pms.ticker)
		pms.ticker = nil
	}
	return err
}

func (pms *PMS) PingConnect() error {
	var err error

	addr := makeAddress(pms.host, pms.port)

	if pms.MpdClient != nil {
		err = pms.MpdClient.Ping()
		if err != nil {
			console.Log("MPD control connection timeout.")
		}
	}

	if pms.MpdClient == nil || err != nil {
		console.Log("Establishing MPD control connection to %s...", addr)
		pms.MpdClient, err = mpd.DialAuthenticated(`tcp`, addr, pms.password)
		if err != nil {
			pms.Error("MPD control connection error: %s", err)
		}
		console.Log("Connected to %s.", addr)
	}

	return err
}

// Monitor connection for errors and terminate when an error occurs
func (pms *PMS) watchMpdIdleErrors() {
	for err := range pms.MpdClientWatcher.Error {
		pms.Error("Error in MPD IDLE connection: %s", err)
		pms.MpdClient.Close()
		pms.MpdClientWatcher.Close()
	}
	go pms.LoopConnect()
}

// Watch for IDLE events and trigger actions when events arrive
func (pms *PMS) watchMpdIdleEvents() {
	var err error

	for subsystem := range pms.MpdClientWatcher.Event {

		console.Log("MPD says it has IDLE events on the following subsystem: %s", subsystem)
		if pms.PingConnect() != nil {
			pms.Error("IDLE: failed to establish MPD control connection: going out of sync with MPD!")
			continue
		}

		switch subsystem {
		case "database":
			err = pms.SyncLibrary()
		case "playlist":
			err = pms.SyncQueue()
		case "player":
			err = pms.UpdatePlayerStatus()
			if err != nil {
				break
			}
			err = pms.UpdateCurrentSong()
		case "options":
			err = pms.UpdatePlayerStatus()
		case "mixer":
			err = pms.UpdatePlayerStatus()
		default:
			console.Log("Ignoring updates by subsystem %s", subsystem)
		}
		if err != nil {
			pms.Error("Error updating status: %s", err)
		}
	}
}

func (pms *PMS) CurrentMpdClient() *mpd.Client {
	err := pms.PingConnect()
	if err == nil {
		return pms.MpdClient
	}
	return nil
}

// runTicker starts a ticker that will increase the elapsed time every second.
func (pms *PMS) runTicker() {
	pms.ticker = make(chan time.Time, 0)

	go func() {
		ticker := time.NewTicker(time.Millisecond * 1000)
		defer ticker.Stop()
		for t := range ticker.C {
			pms.ticker <- t
		}
	}()
	for _ = range pms.ticker {
		pms.MpdStatus.Tick()
		pms.EventPlayer <- 0
	}
}

// Sync retrieves the MPD library and stores it as a Songlist in the
// PMS.Library variable. Furthermore, the search index is opened, and if it is
// older than the database version, a reindex task is started.
//
// If the Songlist or Index is cached at the correct version, that part goes untouched.
func (pms *PMS) SyncLibrary() error {
	if pms.MpdClient == nil {
		return fmt.Errorf("Cannot call Sync() while not connected to MPD")
	}

	stats, err := pms.MpdClient.Stats()
	if err != nil {
		return fmt.Errorf("Error while retrieving library stats from MPD: %s", err)
	}

	libraryVersion, err := strconv.Atoi(stats["db_update"])
	console.Log("Sync(): server reports library version %d", libraryVersion)
	console.Log("Sync(): local version is %d", pms.libraryVersion)

	if libraryVersion != pms.libraryVersion {
		pms.Message("Retrieving library metadata, %s songs...", stats["songs"])
		library, err := pms.retrieveLibrary()
		if err != nil {
			return fmt.Errorf("Error while retrieving library from MPD: %s", err)
		}
		pms.Library = library
		pms.libraryVersion = libraryVersion
		pms.Message("Library metadata at version %d.", pms.libraryVersion)
		pms.EventLibrary <- 1
	}

	console.Log("Sync(): opening search index")
	err = pms.openIndex()
	if err != nil {
		return fmt.Errorf("Error while opening index: %s", err)
	}
	console.Log("Sync(): index at version %d", pms.indexVersion)
	pms.EventIndex <- 1

	if libraryVersion != pms.indexVersion {
		console.Log("Sync(): index version differs from library version, reindexing...")
		err = pms.ReIndex()
		if err != nil {
			return fmt.Errorf("Failed to reindex: %s", err)
		}

		err = pms.writeIndexStateFile(pms.indexVersion)
		if err != nil {
			console.Log("Sync(): couldn't write index state file: %s", err)
		}
		console.Log("Sync(): index updated to version %d", pms.indexVersion)
	}

	console.Log("Sync(): finished.")

	return nil
}

func (pms *PMS) SyncQueue() error {
	if err := pms.UpdatePlayerStatus(); err != nil {
		return err
	}
	if pms.queueVersion == pms.MpdStatus.Playlist {
		return nil
	}
	pms.Message("Retrieving queue metadata, %d songs...", pms.MpdStatus.PlaylistLength)
	queue, err := pms.retrieveQueue()
	if err != nil {
		return fmt.Errorf("Error while retrieving queue from MPD: %s", err)
	}
	pms.Queue = queue
	pms.queueVersion = pms.MpdStatus.Playlist
	pms.Message("Queue at version %d.", pms.queueVersion)
	pms.EventQueue <- 1
	return nil
}

func (pms *PMS) retrieveLibrary() (*songlist.Library, error) {
	timer := time.Now()
	list, err := pms.MpdClient.ListAllInfo("/")
	if err != nil {
		return nil, err
	}
	console.Log("ListAllInfo in %s", time.Since(timer).String())

	s := songlist.NewLibrary()
	s.AddFromAttrlist(list)
	s.SetName("Library")
	return s, nil
}

func (pms *PMS) retrieveQueue() (*songlist.Queue, error) {
	timer := time.Now()
	list, err := pms.MpdClient.PlaylistInfo(-1, -1)
	if err != nil {
		return nil, err
	}
	console.Log("PlaylistInfo in %s", time.Since(timer).String())

	s := songlist.NewQueue()
	s.AddFromAttrlist(list)
	s.SetName("Queue")
	return s, nil
}

func (pms *PMS) openIndex() error {
	timer := time.Now()
	index_dir := indexDirectory(pms.host, pms.port)
	err := createDirectory(index_dir)
	if err != nil {
		return fmt.Errorf("Unable to create index directory %s!", index_dir)
	}

	if pms.Index != nil {
		pms.Index.Close()
	}

	pms.Index, err = index.New(index_dir, pms.Library)
	if err != nil {
		return fmt.Errorf("Unable to acquire index: %s", err)
	}

	pms.indexVersion, err = pms.readIndexStateFile()
	if err != nil {
		console.Log("Sync(): couldn't read index state file: %s", err)
	}

	console.Log("Opened search index in %s", time.Since(timer).String())

	return nil
}

func (pms *PMS) CurrentSong() *song.Song {
	pms.mutex.Lock()
	defer pms.mutex.Unlock()
	return pms.currentSong
}

// UpdateCurrentSong stores a local copy of the currently playing song.
func (pms *PMS) UpdateCurrentSong() error {
	attrs, err := pms.MpdClient.CurrentSong()
	if err != nil {
		return err
	}

	console.Log("MPD current song: %s", attrs["file"])

	pms.mutex.Lock()
	pms.currentSong = song.New()
	pms.currentSong.SetTags(attrs)
	pms.mutex.Unlock()

	pms.EventPlayer <- 0

	return nil
}

// UpdatePlayerStatus populates pms.MpdStatus with data from the MPD server.
func (pms *PMS) UpdatePlayerStatus() error {
	attrs, err := pms.MpdClient.Status()
	if err != nil {
		return err
	}

	pms.MpdStatus.SetTime()

	console.Log("MPD player status: %s", attrs)

	pms.MpdStatus.Audio = attrs["audio"]
	pms.MpdStatus.Err = attrs["err"]
	pms.MpdStatus.State = attrs["state"]

	// The time field is divided into ELAPSED:LENGTH.
	// We only need the length field, since the elapsed field is sent as a
	// floating point value.
	split := strings.Split(attrs["time"], ":")
	if len(split) == 2 {
		pms.MpdStatus.Time, _ = strconv.Atoi(split[1])
	} else {
		pms.MpdStatus.Time = -1
	}

	pms.MpdStatus.Bitrate, _ = strconv.Atoi(attrs["bitrate"])
	pms.MpdStatus.Playlist, _ = strconv.Atoi(attrs["playlist"])
	pms.MpdStatus.PlaylistLength, _ = strconv.Atoi(attrs["playlistlength"])
	pms.MpdStatus.Song, _ = strconv.Atoi(attrs["song"])
	pms.MpdStatus.SongID, _ = strconv.Atoi(attrs["songid"])
	pms.MpdStatus.Volume, _ = strconv.Atoi(attrs["volume"])

	pms.MpdStatus.Elapsed, _ = strconv.ParseFloat(attrs["elapsed"], 64)
	pms.MpdStatus.MixRampDB, _ = strconv.ParseFloat(attrs["mixrampdb"], 64)

	pms.MpdStatus.Consume, _ = strconv.ParseBool(attrs["consume"])
	pms.MpdStatus.Random, _ = strconv.ParseBool(attrs["random"])
	pms.MpdStatus.Repeat, _ = strconv.ParseBool(attrs["repeat"])
	pms.MpdStatus.Single, _ = strconv.ParseBool(attrs["single"])

	pms.EventPlayer <- 0

	return nil
}

func (pms *PMS) ReIndex() error {
	timer := time.Now()
	if err := pms.Index.IndexFull(); err != nil {
		return err
	}
	pms.indexVersion = pms.libraryVersion
	pms.Message("Song library index complete, took %s", time.Since(timer).String())
	return nil
}

// KeyInput receives key input signals, checks the sequencer for key bindings,
// and runs commands if key bindings are found.
func (pms *PMS) KeyInput(ev parser.KeyEvent) {
	matches := pms.Sequencer.KeyInput(ev)
	seqString := pms.Sequencer.String()
	statusText := seqString

	input := pms.Sequencer.Match()
	if !matches || input != nil {
		// Reset statusbar if there is either no match or a complete match.
		statusText = ""
	}

	pms.UI.App.PostFunc(func() {
		pms.UI.Multibar.SetSequenceText(statusText)
		pms.UI.App.Update()
	})

	if input == nil {
		return
	}

	console.Log("Input sequencer matches bind: '%s' -> '%s'", seqString, input.Command)
	pms.UI.EventInputCommand <- input.Command
}

func (pms *PMS) Execute(cmd string) {
	//console.Log("Input command received from Multibar: %s", cmd)
	err := pms.CLI.Execute(cmd)
	if err != nil {
		pms.EventError <- fmt.Sprintf("%s", err)
	}
}
