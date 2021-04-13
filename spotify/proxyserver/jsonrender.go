package spotify_proxyserver

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

type JSONRenderer struct{}

type Response struct {
	Error error  `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
}

func (r *JSONRenderer) Render(w http.ResponseWriter, code int, err error, token *oauth2.Token) {
	var jsontok = make([]byte, 0)
	var erro error
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	jsontok, erro = json.Marshal(token)
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(Response{
		Error: err,
		Token: base64.RawStdEncoding.EncodeToString(jsontok),
	})
}
