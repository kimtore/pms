package main

import (
	"fmt"
	"github.com/ambientsound/pms/config"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
)

func main() {
	exitCode, err := run()
	if exitCode != ExitSuccess {
		log.Error(err)
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
	var token oauth2.Token

	cfg, err := config.Configuration()
	if err != nil {
		flag.Usage()
		return ExitConfiguration, err
	}

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

	handler := &Handler{
		auth:  auth,
		state: state,
		token: make(chan oauth2.Token, 1),
	}
	go http.ListenAndServe("127.0.0.1:59999", handler)

	if len(cfg.Spotify.AccessToken) > 0 && len(cfg.Spotify.RefreshToken) > 0 {
		token = oauth2.Token{
			AccessToken:  cfg.Spotify.AccessToken,
			RefreshToken: cfg.Spotify.RefreshToken,
		}
	} else {
		url := auth.AuthURL(state)
		log.Printf("Please visit this URL to authenticate to Spotify: %s", url)
	}

	err = viper.WriteConfig()
	if err != nil {
		log.Errorf("Unable to write configuration file: %s", err)
		log.Infof("For quick startup next time, please add the following to your configuration file:")
		log.Infof("spotify:")
		log.Infof("  accesstoken: %s", token.AccessToken)
		log.Infof("  refreshtoken: %s", token.RefreshToken)
	}

	env := Environment{
		tokenCallbackHandler: handler,
		signals:              make(chan os.Signal, 1),
	}

	signal.Notify(env.signals, syscall.SIGTERM, syscall.SIGINT)

	return mainloop(env)
}

type Environment struct {
	tokenCallbackHandler *Handler
	signals              chan os.Signal
	client               spotify.Client
}

func mainloop(env Environment) (int, error) {
	for {
		select {
		case token := <-env.tokenCallbackHandler.token:
			log.Infof("Received new Spotify token")
			env.client = env.tokenCallbackHandler.auth.NewClient(&token)
			viper.Set("spotify.accesstoken", token.AccessToken)
			viper.Set("spotify.refreshtoken", token.RefreshToken)

		case sig := <-env.signals:
			log.Infof("Caught signal %d, exiting.", sig)
			return ExitSuccess, nil
		}
	}
}
