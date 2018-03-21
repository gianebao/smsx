package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gianebao/smsx/shorten"
	"github.com/gianebao/smsx/sms"
)

// Message represents the JSON params required for MessageHandler
type Message struct {
	To      string   `json:"to"`
	Message string   `json:"text"`
	URLs    []string `json:"urls,omitempty"`
}

func throwStatus(w http.ResponseWriter, code int, text string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write((Status{
		Status:  code,
		Text:    text,
		Network: 1,
	}).Bytes())
}

// BuildSMSMessage creates the Message object for sms.Send
func BuildSMSMessage(m Message) (int, string, sms.Message) {
	var (
		err        error
		smsMessage sms.Message
		u          string
	)

	smsMessage = sms.Message{
		Template: m.Message,
		Tokens:   []interface{}{},
	}

	for _, uu := range m.URLs {
		if u, err = shorten.Shorten(bitly, uu); err != nil {
			return http.StatusBadRequest, "INVALID_PARAM_URL", smsMessage
		}
		smsMessage.Tokens = append(smsMessage.Tokens, u)
	}

	return http.StatusOK, "OK", smsMessage
}

// MessageHandler handles the request and response in sending sms
func MessageHandler(w http.ResponseWriter, req *http.Request) {
	var (
		body       []byte
		err        error
		m          Message
		smsMessage sms.Message
		s          interface{}
		code       int
		codeText   string
	)

	if body, err = ioutil.ReadAll(req.Body); err != nil {
		throwStatus(w, http.StatusBadRequest, "INVALID_REQUEST_BODY")
		return
	}

	if err = json.Unmarshal(body, &m); err != nil {
		throwStatus(w, http.StatusBadRequest, "INVALID_REQUEST_BODY")
		return
	}

	if code, codeText, smsMessage = BuildSMSMessage(m); http.StatusOK != code {
		throwStatus(w, code, codeText)
		return
	}

	if s, err = sms.Send(nexmo, m.To, smsMessage); err != nil {
		throwStatus(w, http.StatusInternalServerError, "INTERNAL_GATEWAY_ERROR")
		return
	}

	sNexmo := s.(sms.NexmoResponse)

	if sms.NexmoResponseMessageStatusOK == sNexmo.Messages[0].Status {
		throwStatus(w, http.StatusOK, "OK")
		return
	}

	throwStatus(w, http.StatusInternalServerError, "GATEWAY_"+sNexmo.Messages[0].ErrorText)
}
