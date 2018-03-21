package app_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gianebao/smsx/app"
	"github.com/gianebao/smsx/shorten"
	"github.com/gianebao/smsx/sms"
	"github.com/stretchr/testify/assert"
)

func callMessageHandler(pl string) (*http.Response, []byte) {
	req := httptest.NewRequest(
		"POST",
		"https://api.com",
		strings.NewReader(pl))
	w := httptest.NewRecorder()
	app.MessageHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp, body
}

func TestMessageHandler(t *testing.T) {
	os.Setenv("NEXMOAPIKEY", "nexmo123")
	os.Setenv("NEXMOAPISECRET", "nexmo123NEXMO456nexmo123NEXMO456")
	os.Setenv("BITLYUSERNAME", "bitlyuser")
	os.Setenv("BITLYPASSWORD", "P@$sW0rd!")
	app.Init()
	flag.Parse()

	bitlyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "" != r.Header.Get("Authentication") {
			fmt.Fprintln(w, "111111111111111111111111111111111111111")
		} else {
			fmt.Fprintln(w, "http://bit.ly/a234567")
		}
	}))

	defer bitlyServer.Close()

	shorten.BitlyAccessTokenEndpoint = bitlyServer.URL
	shorten.BitlyEndpoint = bitlyServer.URL

	nexmoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"message-count": "1","messages": [{
					"to": "99887711",
					"message-id": "0F0000008BD3AD66",
					"status": "0",
					"remaining-balance": "1.97600000",
					"message-price": "0.02400000",
					"network": "52501"
			}]}`)
	}))

	defer nexmoServer.Close()

	sms.NexmoEndpoint = nexmoServer.URL

	var (
		sCode       int
		sText       string
		jSMSMessage []byte
		smsMessage  sms.Message
	)

	sCode, sText, smsMessage = app.BuildSMSMessage(app.Message{To: "99887711", Message: "test message"})

	assert.Equal(t, http.StatusOK, sCode)
	assert.Equal(t, "OK", sText)
	jSMSMessage, _ = json.Marshal(smsMessage)
	assert.Equal(t, `{"Template":"test message","Tokens":[]}`, string(jSMSMessage))

	resp, body := callMessageHandler(`{"to":"99887711", "url":"https://google.com"}`)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, `{"status":200,"text":"OK","network":1}`, string(body))

	//resp, err := shorten.Shorten(s, "https://google.com")
	//fmt.Println(resp, err)
}
