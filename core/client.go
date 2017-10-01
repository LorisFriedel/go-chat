package core

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

var (
	ClientSuicide    = errors.New("client doesn't want to live anymore")
	ErrWrongPassword = errors.New("wrong password")
	ErrChannel       = errors.New("channel error")
	ErrUnknown       = errors.New("unknown error")
)

type MsgListener func(Message)

type IClient interface {
	Identity() Identity
	AddListener(listener MsgListener)
	CreateConnectChan(name, address string, port int, password string, timeout int) error
	CreateChan(name, address string, port int, password string, timeout int) error
	Connect(name, address string, port int, password string) error
	ConnectKnown(name string) error
	ListKnownChan() func(string) []string
	ListOwnChan() func(string) []string
	CloseChan(name string) error
	Bye() error
	Die() error
	Forget(name string) error
	Me() error
	List() error
	SendMessage(text string) error
}

type Client struct {
	identity   Identity
	currPipe   *Pipe
	currChan   *KnownChan
	knownChans map[string]*KnownChan
	ownChans   map[string]*Channel
	listeners  []MsgListener
}

func NewClient(name string) *Client {
	return &Client{
		identity:   *NewIdentity(name),
		currPipe:   nil,
		currChan:   nil,
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

// TODO function DeleteListener ?

func (c *Client) listen() {
	defer func() {
		c.currChan = nil
		c.currPipe = nil
	}()

	p := c.currPipe
	for msg, err := p.Read(); p.IsOpen(); msg, err = p.Read() {
		// Synchronisation here, a client can't receive more than one message at once and handle them non-concurrently84
		if err != nil {
			log.Errorf("Client.listen: error while reading message from channel (%v)\n", err)
			continue
		}
		c.notify(msg)
	}
	c.notify(*NewMsgSysClient(c.identity, fmt.Sprintf("Disconnected from channel %v", c.currChan)))
}

func (c *Client) notify(msg Message) {
	for _, listener := range c.listeners {
		listener(msg)
	}
}

func (c *Client) CreateConnectChan(name, address string, port int, password string, timeout int) error {
	err := c.CreateChan(name, address, port, password, timeout)
	if err != nil {
		log.Errorf("Client.CreateConnectChan: channel created but won't open (%v,%v,%v)\n", name, address, port)
		return err
	}

	if c.identity.Name != "" {
		err = c.Connect(name, address, port, password)
		if err != nil {
			log.Errorf("Client.CreateConnectChan: client created a channel but can't connect to it (%v)\n", err)
			return err
		}
		log.Infoln("Client.CreateConnectChan: successfully created channel, client is now connected to it")
	}

	return nil
}

func (c *Client) CreateChan(name, address string, port int, password string, timeout int) error {
	channel := NewChannel(address, port, password, timeout)

	if ch, ok := c.ownChans[name]; ok {
		err := ch.Close()
		if err != nil {
			log.Errorf("Client.CreateChan: error while closing channel %v", ch.Addr())
		}
		log.Infof("Client.CreateChan: channel %v will be replaced", ch.Addr())
	}

	err := channel.Open()
	if err != nil {
		log.Errorf("Client.CreateChan: channel created but won't open (%v)\n", channel)
		return err
	}

	c.ownChans[name] = channel

	log.Infoln("Client.CreateChan: successfully created channel")
	return nil
}

func (c *Client) Connect(name, address string, port int, password string) error {
	// Close previous connection
	if c.currPipe != nil && c.currPipe.IsOpen() {
		c.currPipe.Close()
	}

	// Connect to server address
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Errorf("Client.Connect: connection failure on (%s:%d)\n", address, port)
		return err
	}

	// Create pipe to communicate with the channel
	c.currPipe = NewPipe(conn)

	if _, set := c.knownChans[name]; set {
		log.Infof("Client.Connect: replacing channel %s in known channels\n", name)
	} else {
		log.Infof("Client.Connect: adding channel %s to known channels\n", name)
	}

	kChan := newKnownChan(name, address, port, password)
	c.knownChans[name] = kChan
	c.currChan = kChan

	c.send(*NewMsg(c.identity, HELLO))
	msg, err := c.currPipe.Read()
	if err != nil {
		return err
	}

	switch msg.Type {
	case WELCOME_BACK:
		// Empty, we are connected and already authenticated
	case PASSWORD_PLEASE:
		c.currPipe.Write(*NewMsgPassword(c.identity, password))
		answer, err := c.currPipe.Read()
		if err != nil {
			// TODO handle error
		}
		if answer.Type != WELCOME {
			switch answer.Type {
			case WRONG_PASSWORD:
				return ErrWrongPassword
			case ERROR:
				return ErrChannel
			default:
				return ErrUnknown
			}
		}
		// Authentication OK
	case WELCOME:
		// No password, Authentication OK
	}

	go c.listen()
	c.notify(*NewMsgSysClient(c.identity, fmt.Sprintf("Now connected to %v", kChan)))

	return nil
}

func (c *Client) ConnectKnown(name string) error {
	ch, set := c.knownChans[name]
	if !set {
		log.Errorf("Client.ConnectKnown: client tried to connect to an unexisting known channel (%v)\n", name)
		return fmt.Errorf("unknown channel: %s", name)
	}

	return c.Connect(ch.name, ch.address, ch.port, ch.password)
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

func (c *Client) ListOwnChan() func(string) []string {
	return func(input string) []string {
		names := make([]string, 0)
		for name := range c.ownChans {
			names = append(names, name)
		}
		return names
	}
}

func (c *Client) CloseChan(name string) error {
	// TODO bugfix : sometime, the channel owner cannot close the channel
	var ch *Channel
	var set bool

	if ch, set = c.ownChans[name]; !set {
		log.Errorln("Client.CloseChan: try to close a channel he doesn't own")
		return fmt.Errorf("can't close channel %s : you're not the owner", name)
	}

	err := ch.Close()
	if err != nil {
		log.Errorln("Client.CloseChan: error while closing channel")
		return err
	}

	err = c.Forget(name)
	if err != nil {
		log.Errorln("Client.CloseChan: error while forgetting channel")
		return err
	}

	return nil
}

func (c *Client) Bye() error {
	if c.currPipe == nil || !c.currPipe.IsOpen() {
		return errors.New("not connected to any channel")
	}

	log.Infoln("Client.Bye: disconnecting from current channel")
	// We don't tell the channel we are leaving, he will notice himself
	c.notify(*NewMsgSysClient(c.identity, fmt.Sprintf("Goodbye %v", c.currChan))) // TODO useless (or proper bye notification) ?

	return c.currPipe.Close()
}

func (c *Client) Die() error {
	log.Infoln("Client.Die: end of the client instance")
	c.Bye()
	return ClientSuicide
}

func (c *Client) Forget(name string) error {
	_, set := c.knownChans[name]
	if !set {
		log.Errorf("Client.ConnectKnown: client tried to forget unknown channel (%v)\n", name)
		return fmt.Errorf("unknown channel: %s", name)
	}

	delete(c.knownChans, name)
	// TODO Better handle name collision
	return nil
}

func (c *Client) Me() error {
	var text string
	if c.currPipe != nil && c.currPipe.IsOpen() {
		text = fmt.Sprintf("Currently connected to %v", c.currChan)
	} else {
		text = fmt.Sprint("Not connected to any channel :(")
	}
	c.notify(*NewMsgSysClient(c.identity, text))
	return nil
}

func (c *Client) List() error {
	var buffer bytes.Buffer
	buffer.WriteString("List of kown channels:\n")
	for _, ch := range c.knownChans {
		buffer.WriteString(fmt.Sprintln(ch))
	}
	c.notify(*NewMsgText(c.identity, buffer.String()))
	return nil
}

func (c *Client) SendMessage(text string) error {
	log.Infoln("Client.SendMessage: sending message")
	fmt.Printf("\033[1A\033[K")
	return c.send(*NewMsgText(c.identity, text))
}

func (c *Client) send(msg Message) error {
	if c.currPipe == nil || !c.currPipe.IsOpen() {
		log.Errorln("Client.SendMessage: client is not connected to any channel")
		return errors.New("client is not connected to any channel, can't send message")
	}

	return c.currPipe.Write(msg)
}

func (c *Client) String() string {
	return c.identity.Name
}

/*******************/

type KnownChan struct {
	name     string
	address  string
	port     int
	password string
}

func newKnownChan(name, address string, port int, password string) *KnownChan {
	return &KnownChan{
		name:     name,
		address:  address,
		port:     port,
		password: password,
	}
}

func (k *KnownChan) String() string {
	return fmt.Sprintf("%s (%s:%d)", k.name, k.address, k.port)
}
