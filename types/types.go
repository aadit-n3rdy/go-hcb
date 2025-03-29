package types

import "time"

type Message struct {
	Subject   string
	Value     string
	Timestamp time.Time
}

func NewMessage(sub string, val string) *Message {
	msg := new(Message)
	msg.Subject = sub
	msg.Value = val
	msg.Timestamp = time.Now()
	return msg
}
