package core

import (
	"time"
)

type tMsg int

const (
	HELLO           tMsg = 0
	TEXT            tMsg = 1
	WELCOME         tMsg = 2
	WELCOME_BACK    tMsg = 3
	PASSWORD_PLEASE tMsg = 4
	PASSWORD        tMsg = 5
	WRONG_PASSWORD  tMsg = 6
	ERROR           tMsg = 7
	SYS_CLIENT      tMsg = 8
	SYS_CHANNEL     tMsg = 9
)

type Message struct {
	// Must be public for marshalling
	Text      string
	Type      tMsg
	Sender    Identity
	Timestamp time.Time
}

func NewMsgText(sender Identity, text string) *Message {
	return &Message{
		Text:      text,
		Type:      TEXT,
		Sender:    sender,
		Timestamp: time.Now(),
	}
}

func NewMsgSysClient(sender Identity, text string) *Message {
	return &Message{
		Text:      text,
		Type:      SYS_CLIENT,
		Sender:    sender,
		Timestamp: time.Now(),
	}
}

func NewMsgSysChannel(sender Identity, text string) *Message {
	return &Message{
		Text:      text,
		Type:      SYS_CHANNEL,
		Sender:    sender,
		Timestamp: time.Now(),
	}
}

func NewMsgPassword(sender Identity, password string) *Message {
	return &Message{
		Text:      password,
		Type:      PASSWORD,
		Sender:    sender,
		Timestamp: time.Now(),
	}
}

func NewMsg(sender Identity, t tMsg) *Message {
	return &Message{
		Type:      t,
		Sender:    sender,
		Timestamp: time.Now(),
	}
}
