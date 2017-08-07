package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ambientsound/pms/console"
	"github.com/ambientsound/pms/pms"
	"github.com/ambientsound/pms/version"
	"github.com/ambientsound/pms/xdg"

	"github.com/jessevdk/go-flags"
)

var buildVersion = "undefined"

type cliOptions struct {
	Version     bool   `short:"v" long:"version" description:"Print program version"`
	Debug       string `short:"d" long:"debug" description:"Write debugging info to file"`
	MpdHost     string `short:"h" long:"host" description:"MPD host" default-mask:"MPD_HOST environment variable or localhost"`
	MpdPort     string `short:"p" long:"port" description:"MPD port" default-mask:"MPD_PORT environment variable or 6600"`
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

	if len(opts.MpdHost) == 0 {
		val, ok := os.LookupEnv("MPD_HOST")
		if ok {
			opts.MpdHost = val
		} else {
			opts.MpdHost = "localhost"
		}
	}
	if len(opts.MpdPort) == 0 {
		val, ok := os.LookupEnv("MPD_PORT")
		if ok {
			opts.MpdPort = val
		} else {
			opts.MpdPort = "6600"
		}
	}

	pms := pms.New()
	defer func() {
		pms.QuitSignal <- 0
	}()

	// Source default configuration.
	pms.Message("Applying default configuration.")
	if err := pms.SourceDefaultConfig(); err != nil {
		panic(fmt.Sprintf("BUG in default config: %s\n", err))
	}

	// Source configuration files from all XDG standard directories.
	configDirs := xdg.ConfigDirectories()
	for _, dir := range configDirs {
		p := path.Join(dir, "pms.conf")
		pms.Message("Reading configuration file '%s'.", p)
		err = pms.SourceConfigFile(p)
		if err != nil {
			pms.Error("Error while reading configuration file '%s': %s", p, err)
		}
	}

	pms.SetConnectionParams(opts.MpdHost, opts.MpdPort, opts.MpdPassword)
	go pms.LoopConnect()

	pms.Main()
	pms.Wait()

	console.Log("Exiting normally.")
}
