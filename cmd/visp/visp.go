package main

import (
	"fmt"
	"github.com/ambientsound/pms/config"
	"github.com/ambientsound/pms/log"
	"github.com/ambientsound/pms/prog"
	"github.com/ambientsound/pms/widgets"
	"github.com/google/uuid"
	flag "github.com/spf13/pflag"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

var buildVersion = "undefined"

var scopes = []string{
	"playlist-modify-private",
	"playlist-modify-public",
	"playlist-read-collaborative",
	"playlist-read-private",
	"user-follow-modify",
	"user-follow-read",
	"user-library-modify",
	"user-library-read",
	"user-modify-playback-state",
	"user-read-currently-playing",
	"user-read-playback-state",
	"user-read-recently-played",
	"user-top-read",
}

const (
	ExitSuccess = iota
	ExitConfiguration
	ExitInternalError
	ExitPanic
	ExitLogging
)

func logAndStderr(line string) {
	log.Errorf(line)
	fmt.Fprintln(os.Stderr,line)
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

type Handler struct {
	auth  spotify.Authenticator
	state string
	token chan oauth2.Token
}

// the user will eventually be redirected back to your redirect URL
// typically you'll have a handler set up like the following:
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	token, err := h.auth.Token(h.state, r)
	if err != nil || token == nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}

	h.token <- *token
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

	u, err := uuid.NewRandom()
	if err != nil {
		return ExitInternalError, err
	}
	state := u.String()

	auth := spotify.NewAuthenticator("http://localhost:59999/callback", scopes...)
	auth.SetAuthInfo(cfg.Spotify.ClientID, cfg.Spotify.ClientSecret)

	tokens := make(chan oauth2.Token, 1)

	handler := &Handler{
		auth:  auth,
		state: state,
		token: tokens,
	}
	go http.ListenAndServe("127.0.0.1:59999", handler)

	if len(cfg.Spotify.AccessToken) > 0 && len(cfg.Spotify.RefreshToken) > 0 {
		handler.token <- oauth2.Token{
			AccessToken:  cfg.Spotify.AccessToken,
			RefreshToken: cfg.Spotify.RefreshToken,
		}
	} else {
		url := auth.AuthURL(state)
		log.Printf("Please visit this URL to authenticate to Spotify: %s", url)
	}

	visp := &prog.Visp{
		Auth:   auth,
		Tokens: tokens,
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
