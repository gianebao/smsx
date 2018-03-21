package shorten_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gianebao/smsx/shorten"
	"github.com/stretchr/testify/assert"
)

func TestBitly_Shorten(t *testing.T) {
	s := shorten.Bitly{
		Username: "somebitlyuser",
		Password: `QWEUIYQWEIUASGDHA2323`,
	}

	atServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "111111111111111111111111111111111111111")
	}))

	defer atServer.Close()

	shortenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "http://bit.ly/a234567")
	}))

	defer shortenServer.Close()

	shorten.BitlyAccessTokenEndpoint = atServer.URL
	shorten.BitlyEndpoint = shortenServer.URL

	resp, err := shorten.Shorten(s, "https://google.com")
	assert.Equal(t, "http://bit.ly/a234567", resp)
	assert.Nil(t, err)
}
