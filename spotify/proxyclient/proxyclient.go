package spotify_proxyclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ambientsound/pms/spotify/proxyserver"
	"golang.org/x/oauth2"
)

const (
	SubtractTTL = time.Minute * 2
)

// How much time an oauth2 token has to live before it should be refreshed.
func TokenTTL(token *oauth2.Token) time.Duration {
	return token.Expiry.Sub(time.Now()) - SubtractTTL
}

// Decode a base64+JSON encoded oauth2 token into a struct.
func DecodeTokenString(tokenstr string) (*oauth2.Token, error) {
	data, err := base64.RawURLEncoding.DecodeString(tokenstr)
	if err != nil {
		return nil, fmt.Errorf("decode base64 string: %w", err)
	}

	tok := &oauth2.Token{}
	err = json.Unmarshal(data, tok)
	if err != nil {
		return nil, fmt.Errorf("token error: %w", err)
	}

	return tok, nil
}

// Refresh a Spotify oauth2 token using the visp Spotify authorization proxy server.
func RefreshToken(server string, client *http.Client, token *oauth2.Token) (*oauth2.Token, error) {
	// Construct endpoint url
	u, err := url.Parse(server)
	if err != nil {
		return nil, fmt.Errorf("parse Spotify auth url: %s", err)
	}
	u.Path = spotify_proxyserver.RefreshURL

	// Encode request JSON
	payload, err := json.Marshal(token)
	if err != nil {
		return nil, fmt.Errorf("encode token: %s", err)
	}

	// Make the request
	resp, err := client.Post(u.String(), "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("authorization proxy error: %s", err)
	}

	// Read response JSON
	response := &spotify_proxyserver.Response{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, fmt.Errorf("decode data from authorization proxy: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error from authorization proxy: %s: %s", resp.Status, response.Error)
	}

	return DecodeTokenString(response.Token)
}
