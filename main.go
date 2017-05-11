package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/pms"
	"github.com/ambientsound/pms/version"

	"github.com/jessevdk/go-flags"
)

var buildVersion = "undefined"

type cliOptions struct {
	Version     bool   `short:"v" long:"version" description:"Print program version"`
	Debug       string `short:"d" long:"debug" description:"Write debugging info to file"`
	MpdHost     string `long:"host" description:"MPD host (MPD_HOST environment variable)" default:"localhost"`
	MpdPort     string `long:"port" description:"MPD port (MPD_PORT environment variable)" default:"6600"`
	MpdPassword string `long:"password" description:"MPD password"`
}

func main() {
	var opts cliOptions

	version.SetVersion(buildVersion)
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

	if len(opts.Debug) > 0 {
		err := console.Open(opts.Debug)
		if err != nil {
			fmt.Printf("Error while opening log file: %s", err)
			os.Exit(1)
		}
	}

	if opts.Version {
		os.Exit(0)
	}

	console.Log("Starting Practical Music Search.")

	val, ok := os.LookupEnv("MPD_HOST")
	if ok {
		opts.MpdHost = val
	}
	val, ok = os.LookupEnv("MPD_PORT")
	if ok {
		opts.MpdPort = val
	}

	pms := pms.New()
	defer func() {
		pms.QuitSignal <- 0
	}()

	pms.SetConnectionParams(opts.MpdHost, opts.MpdPort, opts.MpdPassword)
	go pms.LoopConnect()

	pms.Main()
	pms.Wait()

	console.Log("Exiting normally.")
}
