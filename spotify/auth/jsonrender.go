package spotify_auth

import (
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

type JSONRenderer struct{}

type Response struct {
	Error error         `json:"error,omitempty"`
	Token *oauth2.Token `json:"token,omitempty"`
}

func (r *JSONRenderer) Render(w http.ResponseWriter, code int, err error, token *oauth2.Token) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Error: err,
		Token: token,
	})
}
