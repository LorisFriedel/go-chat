package core

import (
	"time"
)

type Message struct {
	Text      string
	Sender    Identity
	Timestamp time.Time
}

func newMessage(text string, sender Identity) *Message {
	return &Message{
		Text:      text,
		Sender:    sender,
		Timestamp: time.Now(), // TODO Here of somewhere else ?
	}
}
