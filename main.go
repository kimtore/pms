package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/pms"
	"github.com/ambientsound/pms/version"
	"github.com/ambientsound/pms/widgets"

	"github.com/fhs/gompd/mpd"

	"github.com/jessevdk/go-flags"
)

var build_version string = "undefined"

type Options struct {
	Version     bool   `short:"v" long:"version" description:"Print program version"`
	MpdHost     string `long:"host" description:"MPD host (MPD_HOST environment variable)" default:"localhost"`
	MpdPort     string `long:"port" description:"MPD port (MPD_PORT environment variable)" default:"6600"`
	MpdPassword string `long:"password" description:"MPD password"`
}

func dial(mpd_host, mpd_port string) (c *mpd.Client, err error) {
	addr := fmt.Sprintf("%s:%s", mpd_host, mpd_port)
	c, err = mpd.Dial(`tcp`, addr)
	return
}

func main() {
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

	pms := pms.New()

	timer = time.Now()
	ui := widgets.NewUI()
	ui.Start()
	defer ui.Quit()
	console.Log("UI initialized in %s", time.Since(timer).String())

	go func() {
		err := pms.Connect(opts.MpdHost, opts.MpdPort, opts.MpdPassword)
		if err != nil {
			console.Log("%s", err)
		}
	}()

	go func() {
		for {
			select {
			case <-pms.EventLibrary:
				console.Log("Song library updated in MPD, assigning to UI")
				ui.App.PostFunc(func() {
					ui.Songlist.SetSongList(pms.Library)
					ui.SetDefaultSonglist(pms.Library)
					ui.App.Update()
				})
			case <-pms.EventIndex:
				console.Log("Search index updated, assigning to UI")
				ui.App.PostFunc(func() {
					ui.SetIndex(pms.Index)
				})
			case <-pms.EventPlayer:
				console.Log("Player state has changed")
				ui.App.PostFunc(func() {
					ui.Playbar.SetPlayerStatus(pms.MpdStatus)
					ui.Playbar.SetSong(pms.CurrentSong)
					ui.App.Update()
				})
			}
		}
	}()

	ui.Wait()

	console.Log("Exiting normally.")
}
