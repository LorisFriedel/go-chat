package core

import (
	"fmt"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

// TODO interface ?

type Channel struct {
	open     bool
	wg       sync.WaitGroup
	id       Identity
	registry IRegistry
	address  string
	port     int
	password string
	listener net.Listener
}

func NewChannel(address string, port int, passwd string) *Channel {
	return &Channel{
		open:     false,
		id:       *NewIdentity(fmt.Sprintf("%s:%d", address, port)),
		registry: NewRegistry(),
		address:  address,
		port:     port,
		password: passwd,
	}
}

// Open make the channel listen for connection and handling received message
func (c *Channel) Open() error {
	err := c.listen()
	if err != nil {
		log.Errorln("Channel.Open: can't open channel connection")
		return err
	}

	log.Infof("Channel oppened on %s\n", c.listener.Addr().String())
	return nil
}

func (c *Channel) listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.address, c.port))

	if err != nil {
		log.Errorf("Client.listen: can't listen on (%s)\n", fmt.Sprintf("%s:%d", c.address, c.port))
		return err
	}

	c.listener = listener
	c.open = true

	go c.handleJoin()

	return nil
}

func (c *Channel) handleJoin() {
	c.wg.Add(1)
	for c.open {
		conn, err := c.listener.Accept()
		if err != nil {
			log.Errorf("Channel.listen: connection error: %v", err)
			continue
		}
		go c.Join(conn)
		log.Infof("Channel.listen: %v connected to channel", conn.LocalAddr().String())
	}
	c.wg.Done()
	log.Infoln("Channel.listen: join handling is now inactive")
}

func (c *Channel) Close() (err error) {
	// End infinite loop
	c.open = false

	// Close tcp connection
	err = c.listener.Close()

	// Wait for all loop to properly end
	c.wg.Wait()

	// Close client pipes
	c.registry.Foreach(func(id Identity, p *Pipe) {
		if p.IsOpen() {
			p.Close()
		}
	})

	return
}

// Handle is used to handle message received from connected client
func (c *Channel) Handle(msg Message) {
	// TODO Handle regarding message type ?
	c.Broadcast(msg)
}

// Broadcast send the given message to every client connected to the channel
func (c *Channel) Broadcast(msg Message) {
	log.Infof("Channel.Broadcast: broadcasting message from: %s (%v)", msg.Sender, msg.Timestamp)

	c.registry.Foreach(func(id Identity, p *Pipe) {
		if p.IsOpen() {
			p.Write(msg) // TODO not ignore error ? I mean, who cares ?
		}
	})
}

// Join create a pipe between the channel and the given connection
// and start listening for client message, after executing the authentication procedure
// (the client must know the protocol)
func (c *Channel) Join(conn net.Conn) {
	// Open pipe to communicate with the client
	p := NewPipe(conn)

	// Listen for HELLO message, with client Identity
	msgHello, err := p.Read()
	if err != nil || msgHello.Type != HELLO {
		return
	}

	id := msgHello.Sender

	// TODO Message type useless ?

	// If user is already in our registry, he is authenticated no need for password, send welcome back message
	if c.registry.Exists(id) {
		p.Write(*NewMsg(c.id, WELCOME_BACK)) // TODO handle write error
	} else { // If not, ask for password
		p.Write(*NewMsg(c.id, PASSWORD_PLEASE)) // TODO handle write error
		msgPassword, err := p.Read()
		if err != nil || msgPassword.Type != PASSWORD {
			p.Write(*NewMsg(c.id, ERROR)) // TODO handle write error
			return
		}

		// If password match, OK
		if msgPassword.Text == c.password {
			p.Write(*NewMsg(c.id, WELCOME)) // TODO handle write error
		} else { // If not, close connection.
			p.Write(*NewMsg(c.id, WRONG_PASSWORD)) // TODO handle write error
			return
		}
	}

	c.registry.Push(id, p)
	c.Broadcast(*NewMsgText(c.id, fmt.Sprintf("%s joined the channel.", id.Name)))
	log.Infof("Channel.listen: %s joined the channel (%s)", id.Name, id.Hash)

	go func(id Identity, p *Pipe) {
		// While client is connected
		for msg, err := p.Read(); p.IsOpen(); msg, err = p.Read() {
			if err != nil { // TODO handle error in another way ?
				// TODO log error
				continue
			}
			log.Infof("Channel.listen: received message from: %s (%v)", msg.Sender, msg.Timestamp)
			c.Handle(msg)
		}

		// Here client is disconnected, pipe with him is closed
		c.registry.Pop(id)
		c.Broadcast(*NewMsgText(c.id, fmt.Sprintf("%s leaved the channel.", id.Name)))
		log.Infof("Channel.listen: %s leaved the channel (%s)", id.Name, id.Hash)
	}(id, p)
}

// Addr return the ip address of the channel
func (c *Channel) Addr() net.Addr {
	return c.listener.Addr()
}
