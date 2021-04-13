package spotify_auth

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
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

const (
	cookieName  = "token"
	loginURL    = "/oauth/login"
	callbackURL = "/oauth/callback"
)

type Handler struct {
	auth     spotify.Authenticator
	renderer Renderer
}

type Renderer interface {
	Render(w http.ResponseWriter, code int, err error, token *oauth2.Token)
}

func New(clientID, clientSecret, redirectURL string, renderer Renderer) *Handler {
	authenticator := spotify.NewAuthenticator(redirectURL, scopes...)
	authenticator.SetAuthInfo(clientID, clientSecret)

	return &Handler{
		auth:     authenticator,
		renderer: renderer,
	}
}

func (h *Handler) ServeLogin(w http.ResponseWriter, r *http.Request) {
	u, err := uuid.NewRandom()
	if err != nil {
		h.renderer.Render(w, http.StatusServiceUnavailable, err, nil)
		return
	}

	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   u.String(),
		Expires: time.Now().Add(time.Hour),
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, h.auth.AuthURL(u.String()), http.StatusFound)

	return
}

func (h *Handler) ServeCallback(w http.ResponseWriter, r *http.Request) {
	// Get state parameter from cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		h.renderer.Render(w, http.StatusBadRequest, err, nil)
		return
	}

	// Exchange credentials into Spotify token
	token, err := h.auth.Token(cookie.Value, r)
	if err != nil {
		h.renderer.Render(w, http.StatusForbidden, err, nil)
		return
	}

	// Return token to client
	h.renderer.Render(w, http.StatusOK, nil, token)
}

func Router(handler *Handler) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get(loginURL, handler.ServeLogin)
	router.Get(callbackURL, handler.ServeCallback)

	return router
}
