package main

import (
	"errors"
	"fmt"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/prog"
	"github.com/ambientsound/pms/spotify/auth"
	"github.com/ambientsound/pms/tokencache"
	"github.com/ambientsound/pms/version"
	"github.com/ambientsound/pms/widgets"
	"github.com/ambientsound/pms/xdg"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

const (
	ConfigFileName = "visp.conf"
	TokenFileName  = "token.json"
)

const (
	ExitSuccess = iota
	ExitInternalError
	ExitPanic
)

func logAndStderr(line string) {
	log.Errorf(line)
	fmt.Fprintln(os.Stderr, line)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logAndStderr("*********************************")
			logAndStderr("****** Visp has crashed!!! ******")
			logAndStderr("*********************************")
			logAndStderr("Please report this bug at the Github project and include the following stack trace:")
			stacktrace := strings.Split(string(debug.Stack()), "\n")
			for _, line := range stacktrace {
				logAndStderr(line)
			}
			os.Exit(ExitPanic)
		}
	}()

	exitCode, err := run()
	if exitCode != ExitSuccess {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(exitCode)
}

func run() (int, error) {
	log.Infof("%s %s starting up", version.Program, version.Version)
	log.Infof("This program was compiled on %s", version.BuildDate().String())

	visp := &prog.Visp{
		Auth: spotify_auth.New(spotify_auth.Authenticator()),
	}

	ui, err := widgets.NewApplication(visp)
	if err != nil {
		return ExitInternalError, err
	}

	ui.Init()
	defer ui.Finish()
	go ui.Poll()

	visp.Termui = ui
	visp.Init()

	err = visp.SourceDefaultConfig()
	if err != nil {
		return ExitInternalError, fmt.Errorf("read default configuration: %s", err)
	}

	// Source configuration files from all XDG standard directories.
	for _, dir := range xdg.ConfigDirectories() {
		configFile := filepath.Join(dir, ConfigFileName)

		err = visp.SourceConfigFile(configFile)

		if errors.Is(err, os.ErrNotExist) {
			log.Debugf("Ignoring non-existing configuration file %s", configFile)
		} else if err != nil {
			log.Errorf("Error in configuration file %s: %s", configFile, err)
		}
	}

	// In case a token has been cached on disk, restore it to memory.
	configDirs := xdg.ConfigDirectories()
	tokenFile := filepath.Join(configDirs[len(configDirs)-1], TokenFileName)
	visp.Tokencache = tokencache.New(tokenFile)
	token, err := visp.Tokencache.Read()
	if err == nil && token != nil {
		visp.SetToken(*token)
	} else {
		log.Debugf("Unable to read cached Spotify token: %s", err)
	}

	// Set up a listener for oauth2 authorization code flow.
	go func() {
		err := http.ListenAndServe(spotify_auth.BindAddress, visp.Auth)
		if err != nil {
			panic(err)
		}
	}()

	log.Infof("Ready.")

	err = visp.Main()
	if err != nil {
		return ExitInternalError, err
	}

	return ExitSuccess, nil
}
