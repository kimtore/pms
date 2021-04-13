package tokencache

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
)

type Tokencache struct {
	filename string
	cached   oauth2.Token
}

func New(filename string) Tokencache {
	return Tokencache{filename: filename}
}

func (t *Tokencache) Read() (*oauth2.Token, error) {
	f, err := os.Open(t.filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	token := &oauth2.Token{}
	err = dec.Decode(token)

	t.cached = *token

	return token, err
}

func (t *Tokencache) Write(token oauth2.Token) error {
	if token.AccessToken == t.cached.AccessToken {
		return nil
	}

	f, err := os.OpenFile(t.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(&token)
	if err == nil {
		t.cached = token
	}
	return err
}

func (t *Tokencache) Cached() *oauth2.Token {
	return &t.cached
}
