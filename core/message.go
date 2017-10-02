package core

import (
	"time"
)

type TMsg int

const (
	HELLO           TMsg = 0
	TEXT            TMsg = 1
	WELCOME         TMsg = 2
	WELCOME_BACK    TMsg = 3
	PASSWORD_PLEASE TMsg = 4
	PASSWORD        TMsg = 5
	WRONG_PASSWORD  TMsg = 6
	ERROR           TMsg = 7
	SYS_CLIENT      TMsg = 8
	SYS_CHANNEL     TMsg = 9
)

type Message struct {
	// Must be public for marshalling
	Text      string
	Type      TMsg
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

func NewMsg(sender Identity, t TMsg) *Message {
	return &Message{
		Type:      t,
		Sender:    sender,
		Timestamp: time.Now(),
	}
}
