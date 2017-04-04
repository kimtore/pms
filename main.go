package main

import (
	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/index"
	"github.com/ambientsound/pms/song"
	"github.com/ambientsound/pms/songlist"
	"github.com/ambientsound/pms/version"
	"github.com/ambientsound/pms/widgets"

	"github.com/fhs/gompd/mpd"

	"github.com/jessevdk/go-flags"

	"fmt"
	"os"
	"strings"
	"time"
)

var build_version string = "undefined"

type Options struct {
	Version bool `short:"v" long:"version" description:"Print program version"`
	Index   bool `short:"i" long:"index" description:"Run song re-indexing"`
}

func dial() (c *mpd.Client, err error) {
	mpd_host := os.Getenv("MPD_HOST")
	mpd_port := os.Getenv("MPD_PORT")
	if mpd_host == "" {
		mpd_host = "localhost"
	}
	if mpd_port == "" {
		mpd_port = "6600"
	}
	addr := fmt.Sprintf("%s:%s", mpd_host, mpd_port)
	c, err = mpd.Dial(`tcp`, addr)
	return
}

func main() {
	var songs *songlist.SongList
	var timer time.Time
	var idx index.Index
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

	timer = time.Now()
	ui := widgets.NewUI()
	console.Log("UI initialized in %s", time.Since(timer).String())

	go func() {

		timer = time.Now()
		c, err := dial()
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
		ui.SetDefaultSonglist(songs)

		timer = time.Now()
		path := "/tmp/pms-index.bleve"
		idx, err = index.New(path, songs)
		if err != nil {
			panic(err)
		}
		console.Log("Opened index at %s in %s", path, time.Since(timer).String())
		ui.SetIndex(&idx)

		if opts.Index {
			timer = time.Now()
			idx.IndexFull()
			console.Log("Indexed songlist in %s", time.Since(timer).String())
		}

	}()

	ui.Run()

	console.Log("Exiting normally.")

	/*
		reader := bufio.NewReader(os.Stdin)
		for {
			//fmt.Printf("Query: ")
			text, _ := reader.ReadString('\n')
			if text == "" {
				break
			}
			_, err = idx.Search(text, songs)
			if err != nil {
				//fmt.Printf("Error while searching: %s\n", err)
				continue
			}
		}
	*/

	//fmt.Printf("\n")
}
