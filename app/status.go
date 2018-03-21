package app

import "encoding/json"

// Status represents status response of the app
type Status struct {
	Status  int    `json:"status"`
	Text    string `json:"text"`
	Network int    `json:"network,omitempty"`
}

// Bytes converts Status message into JSON []bytes
func (e Status) Bytes() []byte {
	b, _ := json.Marshal(e)
	return b
}
