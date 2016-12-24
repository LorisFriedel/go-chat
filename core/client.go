package core

import (
	"fmt"
	"github.com/golang/glog"
	"net"
)

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

//func (c *Client) DelListener(listener MsgListener) error {
//	// TODO
//	return nil
//}

func (c *Client) listen() {
	// TODO something to stop it
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

	// TODO !!!!!!!!
	// TODO check if already exist + or already known ???!!
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

	// TODO check if already exist + or already known ???!!
	kChan := newKnownChan(name, channel.Addr().String(), port, passwd)
	c.knownChans[name] = kChan

	err = c.Connect(channel.address, channel.port)
	if err != nil {
		glog.Errorf("Client.CreateChan: client created a channel but can't connect to it (%v)\n", kChan)
		return err
	}

	glog.Infoln("Client.CreateChan: successfully created channel, client is now connected to it")
	return nil
}

func (c *Client) Connect(address string, port int) error {
	// TODO close previous connection
	if c.currPipe != nil && c.currPipe.IsOpen() {

	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		glog.Errorf("Client.Connect: connection failure on (%s:%d)\n", address, port)
		return err
	}

	pipe := newPipe(conn)
	pipe.Open()
	c.currPipe = pipe
	c.listen()

	// TODO : if NOT known, try to connect, if success, add to known and set current, if failure return error
	// TODO check name collisions
	// TODO check if not already know ?
	// TODO try to reach server
	// TODO Encrypt password
	// C'est le serveur qui validera la connexion, a l'envoie des informations de login, si elles sont pas bonne le serveur
	// dira CIAO !!!

	// On se connect, puis on écoute pour une authentification optionnel, si le serveur dit OK tout de suite
	// pas besoin de mot de passe, sinon il y en a besoin

	return nil
}

func (c *Client) ConnectKnown(name string) error {

	// TODO
	return nil
}

func (c *Client) SendMessage(text string) error {
	glog.Infoln("SendMessage: sending message")

	msg := newMessage(text, c.identity)
	// TODO check close channel
	if c.currPipe == nil {
		glog.Errorln("Client.SendMessage: client is not connected to any channel")
		return fmt.Errorf("client is not connected to any channel, can't send message \"%s\"", text)
	}

	c.currPipe.outgoing <- *msg
	// TODO ne pas affichr le message, attendre de le recevoir du serveur
	//fmt.Println(c.currPipe, " : ", msg.text) //TODO !!!!
	return nil
}

// TODO send message [to a channel]

// TODO disconnect [from a channel]
// TODO change current channel, get current satus
// TODO get server status (if owner/or not, info is different)
// TODO change server password (if owner)
// TODO

func (c *Client) String() string {
	return c.identity.Name
}
