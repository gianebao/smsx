package sms

import "fmt"

// Message represents the SMS message
type Message struct {
	Template string
	Tokens   []interface{}
}

func (m Message) String() string {
	return fmt.Sprintf(m.Template, m.Tokens...)
}
