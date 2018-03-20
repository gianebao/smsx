package sms_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gianebao/smsx/sms"
	"github.com/stretchr/testify/assert"
)

func makeNexmoServer(okResp string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		if string(body) == okResp {
			fmt.Fprintln(w, `{"message-count": "1","messages": [{
            "to": "99887711",
            "message-id": "0F0000008BD3AD66",
            "status": "0",
            "remaining-balance": "1.97600000",
            "message-price": "0.02400000",
            "network": "52501"
        }]}`)
		} else {
			fmt.Fprintln(w, `{"message-count": "1",
        "messages": [{
            "status": "4",
            "error-text": "Bad Request"
        }]}`)
		}

	}))
}

func TestNexmo_Send(t *testing.T) {
	n := sms.Nexmo{
		APIKey:    "abcd1234",
		APISecret: "abcd1234WXYZ7890",
		From:      "rdp",
	}

	s := makeNexmoServer(fmt.Sprintf(
		`{"api_key":"%s","api_secret":"%s","from":"%s","to":"99887711","text":"Hello world"}`,
		n.APIKey,
		n.APISecret,
		n.From,
	))
	defer s.Close()

	// Not advisable! Only overwritten for testing
	sms.NexmoEndpoint = s.URL

	response, err := sms.Send(n,
		"99887711",
		sms.Message{
			Template: "Hello %s",
			Tokens:   []interface{}{"world"},
		})

	nResponse := response.(sms.NexmoResponse)

	assert.Nil(t, err)
	assert.Equal(t, "1", nResponse.MessageCount)
	assert.Equal(t, "99887711", nResponse.Messages[0].To)
	assert.Equal(t, "0F0000008BD3AD66", nResponse.Messages[0].MessageID)
	assert.Equal(t, "0", nResponse.Messages[0].Status)
	assert.Equal(t, "1.97600000", nResponse.Messages[0].RemainingBalance)
	assert.Equal(t, "0.02400000", nResponse.Messages[0].MessagePrice)
	assert.Equal(t, "52501", nResponse.Messages[0].Network)
}
