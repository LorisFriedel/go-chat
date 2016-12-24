package core

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"net"
)

var ClientSuicide = errors.New("client doesn't want to live anymore")

type MsgListener func(Message)

// TODO interface ?
type Client struct {
	identity   Identity
	currPipe   *Pipe
	knownChans map[string]*KnownChan
	ownChans   map[string]*Channel
	listeners  []MsgListener
}

func NewClient(name string) *Client {
	return &Client{
		identity:   *newIdentity(name),
		currPipe:   nil,
		knownChans: make(map[string]*KnownChan),
		ownChans:   make(map[string]*Channel),
		listeners:  make([]MsgListener, 0, 1),
	}
}

func (c *Client) Identity() Identity {
	return c.identity
}

func (c *Client) AddListener(listener MsgListener) {
	c.listeners = append(c.listeners, listener)
}

// TODO function DelListener ?

func (c *Client) listen() {
	go func() {
		for msg := range c.currPipe.incoming {
			for _, listener := range c.listeners {
				listener(msg)
			}
		}
	}()
}

func (c *Client) CreateChan(name, address string, port int, passwd string) error {
	channel := newChannel(address, port, passwd)

	if ch, ok := c.ownChans[name]; ok {
		err := ch.Close()
		if err != nil {
			glog.Errorf("Client.CreateChan: error while closing channel %v", ch.Addr())
		}
		glog.Infof("Client.CreateChan: channel %v will be replaced", ch.Addr())
	}

	err := channel.Open()
	if err != nil {
		glog.Errorf("Client.CreateChan: channel created but won't open (%v)\n", channel)
		return err
	}

	c.ownChans[name] = channel

	err = c.Connect(name, channel.address, channel.port, passwd)
	if err != nil {
		glog.Errorf("Client.CreateChan: client created a channel but can't connect to it (%v)\n", err)
		return err
	}

	glog.Infoln("Client.CreateChan: successfully created channel, client is now connected to it")
	return nil
}

func (c *Client) Connect(name, address string, port int, passwd string) error {
	if c.currPipe != nil && c.currPipe.IsOpen() {
		c.currPipe.Close()
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		glog.Errorf("Client.Connect: connection failure on (%s:%d)\n", address, port)
		return err
	}

	// Create pipe to communicate with the channel
	pipe := newPipe(conn)
	pipe.Open()
	c.currPipe = pipe
	c.listen()

	if _, set := c.knownChans[name]; set {
		glog.Infof("Client.Connect: replacing channel %s in known channels\n", name)
	} else {
		glog.Infof("Client.Connect: adding channel %s to known channels\n", name)
	}

	kChan := newKnownChan(name, address, port, passwd)
	c.knownChans[name] = kChan

	// TODO Password
	return nil
}

func (c *Client) ConnectKnown(name string) error {
	ch, set := c.knownChans[name]
	if !set {
		glog.Errorf("Client.ConnectKnown: client tried to connect to an unexisting known channel (%v)\n", name)
		return fmt.Errorf("unknown channel: %s", name)
	}

	return c.Connect(ch.name, ch.address, ch.port, ch.passwd)
}

func (c *Client) ListKnownChan() func(string) []string {
	return func(input string) []string {
		names := make([]string, 0)
		for name := range c.knownChans {
			names = append(names, name)
		}
		return names
	}
}

func (c *Client) Bye() error {
	glog.Infoln("Client.Bye: disconnecting from current channel")
	return c.currPipe.Close()
}

func (c *Client) Die() error {
	glog.Infoln("Client.Die: end of the client instance")
	c.Bye()
	return ClientSuicide
}

func (c *Client) Forget(chanName string) error {
	delete(c.knownChans, chanName)
	return nil
}

func (c *Client) SendMessage(text string) error {
	glog.Infoln("Client.SendMessage: sending message")

	msg := newMessage(text, c.identity)

	if c.currPipe == nil || !c.currPipe.IsOpen() {
		glog.Errorln("Client.SendMessage: client is not connected to any channel")
		return fmt.Errorf("client is not connected to any channel, can't send message \"%s\"", text)
	}

	c.currPipe.outgoing <- *msg
	return nil
}

func (c *Client) String() string {
	return c.identity.Name
}

// TODO get current satus
// TODO get server status (if owner/or not, info is different)
// TODO change server password (if owner)