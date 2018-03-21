package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Nexmo represents a Nexmo request payload (https://developer.nexmo.com)
type Nexmo struct {
	APIKey    string
	APISecret string
	From      string
	to        string
	text      string
}

// NexmoResponse represents a Nexmo response payload in JSON
type NexmoResponse struct {
	MessageCount string                 `json:"message-count,omitempty"`
	Messages     []NexmoResponseMessage `json:"messages,omitempty"`
}

// NexmoResponseMessage represents a Nexmo message within the response payload
type NexmoResponseMessage struct {
	To               string `json:"to,omitempty"`
	MessageID        string `json:"message-id,omitempty"`
	Status           string `json:"status,omitempty"`
	ErrorText        string `json:"error-text,omitempty"`
	RemainingBalance string `json:"remaining-balance,omitempty"`
	MessagePrice     string `json:"message-price,omitempty"`
	Network          string `json:"network,omitempty"`
}

const (
	// NexmoResponseMessageStatusOK defines the Nexmo message status when OK
	NexmoResponseMessageStatusOK = "0"
)

var (
	// NexmoEndpoint defines the Nexmo ReST endpoint
	NexmoEndpoint = "https://rest.nexmo.com/sms/json"
)

// getResponse creates a http server request for Nexmo
func (n Nexmo) getResponse() (NexmoResponse, error) {
	var (
		b      []byte
		err    error
		req    *http.Request
		resp   *http.Response
		nResp  = NexmoResponse{}
		client = &http.Client{}
	)

	// MarshalJSON will not throw any error
	b, _ = json.Marshal(n)

	if req, err = http.NewRequest(http.MethodPost, NexmoEndpoint, bytes.NewReader(b)); err != nil {
		return nResp, err
	}

	if resp, err = client.Do(req); err != nil {
		return nResp, err
	}

	defer resp.Body.Close()

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return nResp, err
	}

	err = json.Unmarshal(b, &nResp)
	return nResp, err
}

// send sends the API request to Nexmo server with the `to` and `message` parameters
func (n Nexmo) send(to string, message Message) (interface{}, error) {

	n.to = to
	n.text = message.String()

	return n.getResponse()
}

// MarshalJSON generates the JSON payload of a Nexmo object
func (n Nexmo) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("{")
	contents := []string{}

	if "" != n.APIKey {
		contents = append(contents, fmt.Sprintf(`"api_key":"%s"`, n.APIKey))
	}

	if "" != n.APISecret {
		contents = append(contents, fmt.Sprintf(`"api_secret":"%s"`, n.APISecret))
	}

	if "" != n.From {
		contents = append(contents, fmt.Sprintf(`"from":"%s"`, n.From))
	}

	if "" != n.to {
		contents = append(contents, fmt.Sprintf(`"to":"%s"`, n.to))
	}

	if "" != n.text {
		contents = append(contents, fmt.Sprintf(`"text":"%s"`, n.text))
	}

	buf.WriteString(strings.Join(contents, ",") + "}")
	return buf.Bytes(), nil
}
