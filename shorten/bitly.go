package shorten

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Bitly represents a Bitly (https://dev.bitly.com) connection object
type Bitly struct {
	Username string
	Password string
}

var (
	// BitlyAccessTokenEndpoint defines the endpoint for requesting access token
	BitlyAccessTokenEndpoint = "https://api-ssl.bitly.com/oauth/access_token"

	// BitlyEndpoint defines the endpoint for requesting shortened URL
	BitlyEndpoint = "https://api-ssl.bitly.com/v3/shorten"

	// AccessToken contains the access token retrieved. Bitly's token never expires
	accessToken = ""
)

func (s Bitly) getAccessToken() (string, error) {
	if accessToken != "" {
		return accessToken, nil
	}

	var (
		b      []byte
		err    error
		req    *http.Request
		resp   *http.Response
		client = &http.Client{}
	)

	if req, err = http.NewRequest(http.MethodPost, BitlyAccessTokenEndpoint, nil); err != nil {
		return "", err
	}

	req.SetBasicAuth(s.Username, s.Password)

	if resp, err = client.Do(req); err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}

	accessToken = string(b)

	return accessToken, err
}

func (s Bitly) shorten(u string) (string, error) {
	var (
		b      []byte
		at     string
		err    error
		req    *http.Request
		resp   *http.Response
		client = &http.Client{}
	)

	if at, err = s.getAccessToken(); err != nil {
		return "", err
	}

	qs := url.Values{}
	qs.Add("access_token", at)
	qs.Add("longUrl", u)
	qs.Add("format", "txt")

	if req, err = http.NewRequest(http.MethodGet, BitlyEndpoint+"?"+qs.Encode(), nil); err != nil {
		return "", err
	}

	if resp, err = client.Do(req); err != nil {
		return "", err
	}

	//resp.Header.Get(key)

	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)

	if http.StatusOK != resp.StatusCode {
		return "", errors.New(strings.TrimSpace(string(b)))
	}

	return strings.TrimSpace(string(b)), err
}
