package spotify_auth

import (
	"github.com/ambientsound/pms/log"
	"github.com/google/uuid"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"net/http"
)

const (
	preState    = "Iaax5Uz/6vIB6cGItlnd/qbDFb/KGcJGmv5XsdD47+vJA6vGznkObqdvb+izbpw1"
	callbackURL = "http://localhost:59999/callback"
	BindAddress = "127.0.0.1:59999"
)

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

type Handler struct {
	auth  spotify.Authenticator
	token chan oauth2.Token
	state string
}

// the user will eventually be redirected back to your redirect URL
// typically you'll have a handler set up like the following:
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	token, err := h.auth.Token(h.state, r)
	if err != nil || token == nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		log.Errorf("Unable to retrieve Spotify token: %s", err)
		return
	}

	w.Write([]byte("Token successfully retrieved. You can close this window now."))
	h.token <- *token
}

func (h *Handler) Tokens() chan oauth2.Token {
	return h.token
}

func (h *Handler) Client(token oauth2.Token) spotify.Client {
	return h.auth.NewClient(&token)
}

func (h *Handler) AuthURL() string {
	h.state = makeState()
	return h.auth.AuthURL(h.state)
}

func Authenticator() spotify.Authenticator {
	return spotify.NewAuthenticator(callbackURL, scopes...)
}

func New(auth spotify.Authenticator) *Handler {
	return &Handler{
		token: make(chan oauth2.Token, 1),
		auth:  auth,
	}
}

func makeState() string {
	u, err := uuid.NewRandom()
	if err != nil {
		return preState
	}
	return u.String()
}
