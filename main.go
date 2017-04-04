package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/version"
	"github.com/ambientsound/pms/widgets"
	"github.com/ambientsound/pms/xdg"

	"github.com/fhs/gompd/mpd"

	"github.com/jessevdk/go-flags"
)

var build_version string = "undefined"

type Options struct {
	Version bool   `short:"v" long:"version" description:"Print program version"`
	Index   bool   `short:"i" long:"index" description:"Run song re-indexing"`
	MpdHost string `long:"host" description:"MPD host and password (MPD_HOST environment variable)" default:"localhost"`
	MpdPort string `long:"port" description:"MPD port (MPD_PORT environment variable)" default:"6600"`
}

func dial(mpd_host, mpd_port string) (c *mpd.Client, err error) {
	addr := fmt.Sprintf("%s:%s", mpd_host, mpd_port)
	c, err = mpd.Dial(`tcp`, addr)
	return
}

func createDirectory(dir string) error {
	dir_mode := os.ModeDir | 0755
	return os.MkdirAll(dir, dir_mode)
}

func openIndex(mpd_host, mpd_port string, library *songlist.SongList) (*index.Index, error) {
	cache_dir := xdg.CacheDirectory()
	index_dir := path.Join(cache_dir, mpd_host, mpd_port, "index")
	err := createDirectory(index_dir)
	if err != nil {
		return nil, fmt.Errorf("Unable to create index directory %s!", index_dir)
	}

	idx, err := index.New(index_dir, library)
	if err != nil {
		return nil, err
	}

	return &idx, nil
}

func main() {
	var songs *songlist.SongList
	var timer time.Time
	var opts Options

	version.SetVersion(build_version)
	fmt.Printf("%s %s\n", version.LongName(), version.Version())

	remainder, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	if len(remainder) > 0 {
		trailing := strings.Join(remainder, " ")
		fmt.Printf("error: trailing characters: %s\n", trailing)
		os.Exit(1)
	}

	if opts.Version {
		os.Exit(0)
	}

	val, ok := os.LookupEnv("MPD_HOST")
	if ok {
		opts.MpdHost = val
	}
	val, ok = os.LookupEnv("MPD_PORT")
	if ok {
		opts.MpdPort = val
	}

	timer = time.Now()
	ui := widgets.NewUI()
	ui.Start()
	defer ui.Quit()
	console.Log("UI initialized in %s", time.Since(timer).String())

	go func() {

		timer = time.Now()
		c, err := dial(opts.MpdHost, opts.MpdPort)
		defer c.Close()
		if err != nil {
			panic(err)
		}
		console.Log("MPD connection in %s", time.Since(timer).String())

		timer = time.Now()
		allsongs, err := c.ListAllInfo("/")
		if err != nil {
			panic(err)
		}
		console.Log("ListAllInfo in %s", time.Since(timer).String())

		timer = time.Now()
		songs = songlist.New()
		for _, attrs := range allsongs {
			s := &song.Song{}
			s.Tags = attrs
			songs.Add(s)
		}
		console.Log("Built library songlist in %s", time.Since(timer).String())
		ui.Songlist.SetSongList(songs)

		timer = time.Now()
		sorted := *songs
		sorted.Sort()
		console.Log("Sorted songlist in %s", time.Since(timer).String())
		ui.Songlist.SetSongList(&sorted)
		ui.SetDefaultSonglist(&sorted)

		timer := time.Now()
		idx, err := openIndex(opts.MpdHost, opts.MpdPort, &sorted)
		if err != nil {
			panic(fmt.Sprintf("Unable to acquire index: %s", err))
		}
		ui.SetIndex(idx)
		console.Log("Opened index for %s:%s in %s", opts.MpdHost, opts.MpdPort, time.Since(timer).String())

		if opts.Index {
			timer = time.Now()
			idx.IndexFull()
			console.Log("Indexed songlist in %s", time.Since(timer).String())
		}

	}()

	ui.Wait()

	console.Log("Exiting normally.")
}
