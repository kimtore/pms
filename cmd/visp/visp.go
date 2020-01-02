package main

import (
	"fmt"
	"github.com/ambientsound/pms/config"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/prog"
	"github.com/ambientsound/pms/spotify/auth"
	"github.com/ambientsound/pms/widgets"
	flag "github.com/spf13/pflag"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var buildVersion = "undefined"

const (
	ExitSuccess = iota
	ExitConfiguration
	ExitInternalError
	ExitPanic
	ExitLogging
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
	cfg, err := config.Configuration()
	if err != nil {
		flag.Usage()
		return ExitConfiguration, err
	}

	err = log.Configure(cfg.Log)
	if err != nil {
		return ExitLogging, fmt.Errorf("error in configuration: %s", err)
	}

	log.Infof("Visp starting up")

	if len(cfg.Spotify.ClientID) == 0 || len(cfg.Spotify.ClientSecret) == 0 {
		flag.Usage()
		return ExitConfiguration, fmt.Errorf("you must configure spotify.clientid and spotify.clientsecret")
	}

	auth := spotify_auth.Authenticator()
	auth.SetAuthInfo(cfg.Spotify.ClientID, cfg.Spotify.ClientSecret)

	handler := spotify_auth.New(auth)
	go func() {
		err := http.ListenAndServe(spotify_auth.BindAddress, handler)
		if err != nil {
			panic(err)
		}
	}()

	if len(cfg.Spotify.AccessToken) > 0 && len(cfg.Spotify.RefreshToken) > 0 {
		handler.Tokens() <- oauth2.Token{
			AccessToken:  cfg.Spotify.AccessToken,
			RefreshToken: cfg.Spotify.RefreshToken,
			Expiry:       time.Now(),
			TokenType:    "Bearer",
		}
	}

	visp := &prog.Visp{
		Auth: handler,
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

	log.Infof("Ready.")

	err = visp.Main()
	if err != nil {
		return ExitInternalError, err
	}

	return ExitSuccess, nil
}
